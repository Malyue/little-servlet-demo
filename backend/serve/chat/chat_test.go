package chat

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"os"
	"testing"
	"time"
)

var index = 1

func CreateWs(num int) {
	url := "ws://127.0.0.1:9090/api/testChat"
	c_url := url + "?userid=" + fmt.Sprint(num)
	conn, _, err := websocket.DefaultDialer.Dial(c_url, nil)
	if err != nil {
		fmt.Println("连接失败")
		fmt.Println(index)
		os.Exit(0)
	}
	fmt.Println("create " + fmt.Sprint(num))

	//消息结构
	message := Message{
		Type:    1,
		Content: "",
	}
	message.Receiver = make([]int, 0)
	message.UserId = make([]int, 0)
	var receiver = 0
	if num%2 == 0 {
		receiver = num - 1
	} else {
		receiver = num + 1
	}
	message.Receiver = append(message.Receiver, receiver)
	message.UserId = append(message.UserId, num)

	go TickerSend(conn, message)
}

// 模拟前端，建立一组ws连接
func CreateChatService(num int) {
	url := "ws://127.0.0.1:9090/api/testChat"

	c_url := url + "?userid=" + fmt.Sprint(num)
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(c_url, nil)
	if err != nil {
		fmt.Println(num, "连接失败", err)
		fmt.Println(2*index-1, "个进程")
		os.Exit(1)
		return
	}
	defer conn.Close()

	c_url = url + "?userid=" + fmt.Sprint(num+1)
	conn2, _, err := websocket.DefaultDialer.Dial(c_url, nil)
	if err != nil {
		fmt.Println(num+1, "连接失败", err)
		fmt.Println(2*index, "进程")
		os.Exit(1)
		return
	}
	defer conn2.Close()

	//发送消息
	message := Message{Type: 1, Content: "Hello"}
	message.Receiver = make([]int, 0)
	message.UserId = make([]int, 0)
	message.Receiver = append(message.Receiver, num+1)
	message.UserId = append(message.UserId, num)

	message2 := Message{Type: 1, Content: "Helo"}
	message2.Receiver = make([]int, 0)
	message2.UserId = make([]int, 0)
	message2.Receiver = append(message2.Receiver, num)
	message2.UserId = append(message2.UserId, num+1)

	index = index + 2

	go TickerSend(conn, message)
	go TickerSend(conn2, message2)

	for {

	}
	//模拟前端发送消息
	//go sendMessage(conn, message)
	//go sendMessage(conn2, message2)

}

func TickerSend(conn *websocket.Conn, message Message) {
	ticker := time.NewTicker(8 * time.Second)
	for {
		<-ticker.C
		message.SendTime = time.Now().Unix()
		byteMessage, _ := json.Marshal(&message)
		err := conn.WriteMessage(websocket.TextMessage, byteMessage)
		if err != nil {
			continue
		}
	}
}

// 循环发送消息
func sendMessage(conn *websocket.Conn, message Message) {
	for {
		message.SendTime = time.Now().Unix()
		byteMessage, _ := json.Marshal(&message)
		err := conn.WriteMessage(websocket.TextMessage, byteMessage)
		if err != nil {
			continue
		}
		time.Sleep(time.Duration(2) * time.Second)
	}
}

//func TestWebSocketClient(t *testing.T) {
//	fmt.Println("begin test")
//
//	//go Listent()
//	//开启一组聊天
//	for {
//		//建立连接
//		go CreateChatService(index)
//		//go CreateChatService(index + 1)
//		//index = index + 2
//		//fmt.Println("开启到第n个:", index-2)
//		time.Sleep(time.Duration(1) * time.Second)
//	}
//}

//func BenchmarkChat(b *testing.B) {
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			conn := CreateWs()
//			conn2 := CreateWs()
//			//发送消息
//			SendMessage(conn, conn2, index-2)
//			time.Sleep(time.Duration(1) * time.Second)
//		}
//	})
//}
//

//func BenchmarkChat(b *testing.B) {
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			Create
//			index = index + 1
//			time.Sleep(time.Duration(1/4) * time.Second)
//		}
//	})
//}

func CreateWsT(num int) {
	url := "ws://127.0.0.1:9090/api/testChat"
	c_url := url + "?userid=" + fmt.Sprint(num)
	dialer := websocket.Dialer{}
	//conn, _, err := websocket.DefaultDialer.Dial(c_url, nil)
	conn, _, err := dialer.Dial(c_url, nil)
	if err != nil {
		fmt.Println("连接失败")
		fmt.Println(index)
		os.Exit(0)
	}
	fmt.Println("create " + fmt.Sprint(num))
	defer conn.Close()
	c_url2 := url + "?userid=" + fmt.Sprint(num+1)
	//conn2, _, err := websocket.DefaultDialer.Dial(c_url2, nil)
	conn2, _, err := dialer.Dial(c_url2, nil)
	defer conn2.Close()
	//消息结构
	message := Message{
		Type:    1,
		Content: "",
	}
	message.Receiver = make([]int, 0)
	message.UserId = make([]int, 0)
	//var receiver = 0
	//if num%2 == 0 {
	//	receiver = num - 1
	//} else {
	//	receiver = num + 1
	//}
	message.Receiver = append(message.Receiver, num+1)
	message.UserId = append(message.UserId, num)

	message2 := Message{
		Type:    1,
		Content: "",
	}
	message2.Receiver = make([]int, 0)
	message2.UserId = make([]int, 0)
	message2.Receiver = append(message.Receiver, num)
	message2.UserId = append(message2.UserId, num+1)
	go TickerSend(conn, message)
	go TickerSend(conn2, message2)
}

func TestChat(t *testing.T) {
	fmt.Println("begin test")
	for {
		CreateWs(index)
		index = index + 1
		time.Sleep(time.Second * 1 / 70)
	}
}
