package chat

import (
	"backend/cmn"
	"fmt"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub    *Hub
	Userid int `json:"userid"`

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		SendUnRegisterInfo(c)
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			fmt.Println(err)
			break
		}
		//解析消息
		//myname, err := GetUserById(c.Userid)
		//if err != nil {
		//	zap.L().Error("获得用户名失败", zap.Error(err))
		//	return
		//}
		sendTime, receiver, content, err := ParseMessage(message)
		fmt.Println("接收时间", time.Now().Sub(sendTime))
		fmt.Println("连接数", len(c.hub.clients))
		if time.Now().Sub(sendTime) > 3*time.Second {
			//fmt.Println("userid", c.Userid)
			fmt.Println("停止")
			os.Exit(0)
		}

		if err != nil {
			zap.L().Error("解析消息错误")
			return
		}
		mymessage := Message{Type: 1, Content: content}
		mymessage.Receiver = make([]int, 0)
		mymessage.UserName = make([]string, 0)
		mymessage.UserId = make([]int, 0)
		mymessage.Receiver = append(mymessage.Receiver, receiver)
		mymessage.UserId = append(mymessage.UserId, c.Userid)
		//mymessage.UserName = append(mymessage.UserName, myname)

		c.hub.Broadcast <- mymessage
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, q *cmn.ServiceCtx) {

	conn, err := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	}}).Upgrade(q.W, q.R, nil)
	//fmt.Println(q.Userid, "升级成功")

	if err != nil {
		fmt.Println(err)
		return
	}
	userid, _ := strconv.Atoi(q.Userid)
	client := &Client{Userid: userid, hub: hub, conn: conn, send: make(chan []byte, 10000)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()

	GetOnlineUser(client)
}

func GetOnlineUser(client *Client) {
	client.hub.ClientsLock.Lock()
	defer client.hub.ClientsLock.Unlock()
	onlineUser := make([]int, 0)
	onlineUserName := make([]string, 0)
	//获得在线人员
	for userid, _ := range client.hub.clients {
		if userid != client.Userid {
			name, err := GetUserById(userid)
			if err != nil {
				zap.L().Error("获得用户名失败", zap.Error(err))
				//fmt.Println("获得用户名失败", err)
				return
			}
			onlineUser = append(onlineUser, userid)
			onlineUserName = append(onlineUserName, name)
		}
	}

	//把自身上线消息推送给其余在线人员
	myname, err := GetUserById(client.Userid)
	if err != nil {
		zap.L().Error("获得用户名失败", zap.Error(err))
		return
	}
	//发给所有在线用户--删除自身
	mymessage := Message{Type: 2}
	mymessage.UserId = make([]int, 0)
	mymessage.UserName = make([]string, 0)
	mymessage.Receiver = make([]int, 0)
	mymessage.UserId = append(mymessage.UserId, client.Userid)
	mymessage.UserName = append(mymessage.UserName, myname)
	mymessage.Receiver = onlineUser

	client.hub.Broadcast <- mymessage

	receiver := make([]int, 0)
	receiver = append(receiver, client.Userid)
	message := Message{Type: 2, UserId: onlineUser, UserName: onlineUserName, Receiver: receiver}
	client.hub.Broadcast <- message
}

func SendUnRegisterInfo(client *Client) {
	client.hub.ClientsLock.Lock()

	//广播自身下线消息
	onlineUser := make([]int, 0)
	for userid, _ := range client.hub.clients {
		onlineUser = append(onlineUser, userid)
	}
	mymessage := Message{Type: 3}
	mymessage.UserId = make([]int, 0)
	mymessage.Receiver = make([]int, 0)
	//fmt.Println(client.Userid)
	mymessage.UserId = append(mymessage.UserId, client.Userid)
	mymessage.Receiver = onlineUser

	client.hub.Broadcast <- mymessage
	defer client.hub.ClientsLock.Unlock()
}
