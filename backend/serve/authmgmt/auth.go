package authmgmt

import (
	"backend/cmn"
	"backend/pkg"
	"backend/tool"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"strconv"
)

//author :{"name":"auth","email":"111111@qq.com"}
//annotation:auth-mgmt-service

//func init() {
//	//z = cmn.GetLogger()
//	zap.L().Error("init auth mgmt service")
//}

func Enroll(author string) {
	zap.L().Info("auth.Enroll called")
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
		Fn: auth,

		Path: "/auth",
		Name: "auth",

		Developer: developer,
	})

	cmn.AddService(&cmn.ServeEndPoint{
		Fn:        login,
		Path:      "/login",
		Name:      "login",
		Developer: developer,
	})
}

func auth(ctx context.Context) {
	fmt.Println("auth")
}

func login(ctx context.Context) {
	q := cmn.GetCtxValue(ctx)

	jsonMap, err := tool.ReadBody(q.R)
	if err != nil {
		zap.L().Error("parse body error", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = cmn.CodeServerBusy.Msg()
		q.Resp()
		return
	}

	account := fmt.Sprint(jsonMap["account"])
	password := fmt.Sprint(jsonMap["password"])

	fmt.Println(account, password)

	id, userPassword, err := CheckIfUser(account)
	if err != nil {
		zap.L().Error("select the pg err", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = cmn.CodeServerBusy.Msg()
		q.Resp()
		return
	}

	if password == userPassword {
		var userClaims pkg.Claims
		userClaims.Id = strconv.Itoa(id)
		token, err := pkg.GenerateToken(&userClaims)
		if err != nil {
			zap.L().Error("generate token err ", zap.Error(err))
			q.Msg.Status = cmn.CodeServerBusy
			q.Msg.Msg = q.Msg.Status.Msg()
			q.Resp()
			return
		}
		q.Msg.Status = cmn.CodeSuccess
		q.Msg.Msg = q.Msg.Status.Msg()
		q.Msg.Data = token
		q.Resp()
		return
	} else {
		q.Msg.Status = cmn.CodeErrLogin
		q.Msg.Msg = cmn.CodeErrLogin.Msg()
		q.Resp()
		return
	}
}
