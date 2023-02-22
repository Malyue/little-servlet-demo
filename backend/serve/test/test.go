package test

//author :{"name":"fileService","email":"1637901557@qq.com"}
//annotation:file-mgmt-service

import (
	"backend/cmn"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
)

func Enroll(author string) {
	zap.L().Info("Enroll.Enroll called")
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
		Fn:   test,
		Path: "/test",
		Name: "testService",

		Developer: developer,
	})
}

func test(ctx context.Context) {
	q := cmn.GetCtxValue(ctx)

	q.Msg.Status = cmn.CodeSuccess
	q.Msg.Msg = cmn.CodeSuccess.Msg()
	q.Resp()
}
