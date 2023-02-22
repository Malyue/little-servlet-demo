package user

import (
	"backend/dao"
	"context"
	"go.uber.org/zap"
	"strings"
	"time"
)

type userInfo struct {
	Id      int       `json:"id"`
	Name    string    `json:"name"`
	Account string    `json:"account"`
	Role    string    `json:"role"`
	Note    string    `json:"note"`
	Time    time.Time `json:"time"`
}

// 返回一定数量的用户信息
func SelectUser(search string, page int, count int) (userinfo []userInfo, err error) {

	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return
	}
	defer dao.Close(conn)

	s := `select uid,username,account,role,createtime from t_user WHERE username like $1 limit $2 offset $3`
	rows, err := conn.Query(context.Background(), s, search, count, (page-1)*count)
	defer rows.Close()
	if err != nil {
		return
	}
	var id int
	var name string
	var account string
	var role string
	var createTime time.Time
	for rows.Next() {
		err = rows.Scan(&id, &name, &account, &role, &createTime)
		role = strings.TrimSpace(role)
		name = strings.TrimSpace(name)
		account = strings.TrimSpace(account)
		if err != nil {
			return
		}
		if role == "1" {
			role = "管理员"
		} else {
			role = "用户"
		}
		info := userInfo{
			Id:      id,
			Time:    createTime,
			Name:    name,
			Role:    role,
			Account: account,
		}
		userinfo = append(userinfo, info)
	}
	return
}

// 返回总用户数量
func SelectCountOfUser(search string) (count int, err error) {
	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return
	}
	defer dao.Close(conn)

	s := `select count(*) from t_user where username like $1`
	rc := conn.QueryRow(context.Background(), s, search)
	err = rc.Scan(&count)

	return
}

// 删除用户
func DeleteUser(userid float64) error {
	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return err
	}
	defer dao.Close(conn)

	s := `delete from t_user where uid = $1`
	_, err = conn.Exec(context.Background(), s, userid)
	return err
}

// 增加用户
func AddUser(account string, password string, username string) error {
	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return err
	}
	defer dao.Close(conn)

	s := `insert into t_user(username,account,password,role) values($1,$2,$3,$4)`
	_, err = conn.Exec(context.Background(), s, username, account, password, "2")
	return err
}

//func Register(password string, userName string, account string) {
//	conn, err := dao.Connect()
//	if err != nil {
//		zap.L().Error("获得数据库连接失败", zap.Error(err))
//		return
//	}
//	defer dao.Close(conn)
//
//	s := `insert into `
//}

//func UpdateUser()

//func SearchUser(search string) (userinfo []userInfo, err error) {
//	conn, err := dao.Connect()
//	if err != nil {
//		zap.L().Error("获得数据库连接失败", zap.Error(err))
//		return
//	}
//	defer dao.Close(conn)
//
//	s := `select`
//}
