package chat

import (
	"encoding/json"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
//var index = 1

type Hub struct {
	// Registered clients.
	clients map[int]*Client

	// Inbound messages from the clients.
	Broadcast chan Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	ClientsLock sync.RWMutex //客户端读写锁
}

type Message struct {
	Type     int      `json:"type"`
	UserId   []int    `json:"userid"`
	UserName []string `json:"userName"`
	Content  string   `json:"content"`
	Receiver []int    `json:"receiver"`
	SendTime int64    `json:"sendTime"`
}

func NewHub() *Hub {
	return &Hub{
		Broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			//fmt.Println(client.Userid, "上线")
			h.ClientsLock.Lock()
			h.clients[client.Userid] = client
			h.ClientsLock.Unlock()
		case client := <-h.unregister:
			h.ClientsLock.Lock()
			//fmt.Println(client.Userid, "下綫")
			if _, ok := h.clients[client.Userid]; ok {
				delete(h.clients, client.Userid)
				close(client.send)
			}
			//fmt.Println(client.Userid, "下线")
			h.ClientsLock.Unlock()
		case message := <-h.Broadcast:
			//fmt.Println(message)
			byteMessage, _ := json.Marshal(&message)
			for i := 0; i < len(message.Receiver); i++ {
				//判断收件人是否在线,在线则直接传输过去
				userClient, ifExist := h.clients[message.Receiver[i]]
				if ifExist {
					userClient.send <- byteMessage
				}
			}
		}
	}
}

func (h *Hub) EventRegister(client *Client) {
	//fmt.Println(client.Userid, "上线")
	//h.ClientsLock.Lock()
	h.clients[client.Userid] = client
	//defer h.ClientsLock.Unlock()
}
