package chat

import (
	"backend/cmn"
	//"backend/service"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

//author :{"name":"chat","email":"1637901557@qq.com"}
//annotation:chat-mgmt-service

func init() {
	//HubManage = NewHub()
	//go HubManage.Run()
}

func Enroll(author string) {
	zap.L().Info("chat.Enroll called")
	var developer *cmn.ModuleAuthor
	if author != "" {
		var d cmn.ModuleAuthor
		err := json.Unmarshal([]byte(author), &d)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		developer = &d
	}

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: chat,

		Path: "/chat",
		Name: "chat",

		Developer:   developer,
		IsWebSocket: true,
	})

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: chat,

		Path:      "/testChat",
		Name:      "TestChat",
		Developer: developer,

		IsWebSocket: true,
	})

}

func chat(ctx context.Context) {
	fmt.Println(ctx)
}

func Chat(ctx context.Context, hub *Hub) {
	q := cmn.GetCtxValue(ctx)
	ServeWs(hub, q)
}
