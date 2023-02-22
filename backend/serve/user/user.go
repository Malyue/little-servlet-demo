package user

import (
	"backend/cmn"
	"backend/serve/authmgmt"
	"backend/tool"
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"reflect"
	"strconv"
)

//author :{"name":"user","email":"1637901557@qq.com"}
//annotation:user-mgmt-service

//func init() {
//	//HubManage = NewHub()
//	//go HubManage.Run()
//}

func Enroll(author string) {
	zap.L().Info("user.Enroll called")
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
		Fn: getUser,

		Path: "/getUser",
		Name: "getUser",

		Developer: developer,
	})

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: deleteUser,

		Path: "/deleteUser",
		Name: "deleteUser",

		Developer: developer,
	})

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: register,

		Path: "/register",
		Name: "register",

		Developer: developer,
	})
	//cmn.AddService(&cmn.ServeEndPoint{
	//	Fn: searchUser,
	//
	//	Path: "/searchUser",
	//	Name: "searchUser",
	//
	//	Developer: developer,
	//})

}

func register(ctx context.Context) {
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
	userName := fmt.Sprint(jsonMap["userName"])

	id, password, err := authmgmt.CheckIfUser(account)
	//如果已存在
	if id != -1 {
		q.Msg.Status = cmn.CodeHadUser
		q.Msg.Msg = cmn.CodeHadUser.Msg()
		q.Resp()
		return
	}

	password = fmt.Sprint(jsonMap["password"])

	err = AddUser(account, password, userName)
	if err != nil {
		zap.L().Error("add user error", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = cmn.CodeServerBusy.Msg()
		q.Resp()
		return
	}

	q.Msg.Status = cmn.CodeSuccess
	q.Msg.Msg = cmn.CodeSuccess.Msg()
	q.Resp()
	return
}

func getUser(ctx context.Context) {
	q := cmn.GetCtxValue(ctx)

	jsonMap, err := tool.ReadBody(q.R)
	if err != nil {
		zap.L().Error("parse body error", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = cmn.CodeServerBusy.Msg()
		q.Resp()
		return
	}

	page, _ := strconv.Atoi(fmt.Sprint(jsonMap["page"]))
	count, _ := strconv.Atoi(fmt.Sprint(jsonMap["count"]))
	searchInput := jsonMap["search"]
	var search string
	if searchInput == nil {
		search = "%"
	} else {
		search = "%" + fmt.Sprint(searchInput) + "%"
	}

	allCount, err := SelectCountOfUser(search)
	if err != nil {
		zap.L().Error("获得所有用户条数失败", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = q.Msg.Status.Msg()
		q.Resp()
		return
	}

	userinfo, err := SelectUser(search, page, count)
	if err != nil {
		zap.L().Error("获得用户信息失败", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = q.Msg.Status.Msg()
		q.Resp()
		return
	}

	type respStruct struct {
		Count    int        `json:"count"`
		UserInfo []userInfo `json:"userInfo"`
	}
	resp := respStruct{
		Count:    allCount,
		UserInfo: userinfo,
	}
	q.Msg.Data = resp
	q.Resp()
	return
}

func deleteUser(ctx context.Context) {
	q := cmn.GetCtxValue(ctx)

	jsonMap, err := tool.ReadBody(q.R)
	if err != nil {
		zap.L().Error("parse body error", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = cmn.CodeServerBusy.Msg()
		q.Resp()
		return
	}

	deleteUserList := jsonMap["deleteUserList"]
	value := reflect.ValueOf(deleteUserList)

	//如果不为数组，则参数错误
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		zap.L().Info("参数错误(需要数组参数)", zap.Error(err))
		q.Msg.Status = cmn.CodeInvalidParam
		q.Msg.Msg = cmn.CodeInvalidParam.Msg()
		q.Resp()
		return
	}

	for i := 0; i < value.Len(); i++ {
		//userid, _ := strconv.Atoi(value.Index(i).Interface().(float64))
		err = DeleteUser(value.Index(i).Interface().(float64))
		if err != nil {
			zap.L().Error("删除用户表失败", zap.Error(err))
			q.Msg.Status = cmn.CodeServerBusy
			q.Msg.Msg = cmn.CodeServerBusy.Msg()
			q.Resp()
			return
		}
	}
	q.Msg.Status = cmn.CodeSuccess
	q.Msg.Msg = cmn.CodeSuccess.Msg()
	q.Resp()
	return
}

//func searchUser(ctx context.Context) {
//	q := cmn.GetCtxValue(ctx)
//
//	jsonMap, err := tool.ReadBody(q.R)
//	if err != nil {
//		zap.L().Error("parse body error", zap.Error(err))
//		q.Msg.Status = cmn.CodeServerBusy
//		q.Msg.Msg = cmn.CodeServerBusy.Msg()
//		q.Resp()
//		return
//	}
//
//	search := fmt.Sprint(jsonMap["search"])
//}
