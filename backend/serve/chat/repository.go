package chat

import (
	"backend/dao"
	"context"
	"go.uber.org/zap"
)

func GetUserById(userid int) (userName string, err error) {
	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return
	}
	defer dao.Close(conn)

	s := `select username from t_user WHERE uid = $1`
	rc := conn.QueryRow(context.Background(), s, userid)
	err = rc.Scan(&userName)

	return userName, err
}
