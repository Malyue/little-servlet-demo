package authmgmt

import (
	"backend/dao"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"strings"
)

func CheckIfUser(account string) (id int, password string, err error) {
	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return
	}
	defer dao.Close(conn)

	s := `select uid,password from t_user WHERE account = $1`
	rc := conn.QueryRow(context.Background(), s, account)
	err = rc.Scan(&id, &password)
	password = strings.TrimSpace(password)

	//如果不存在
	if errors.Is(err, pgx.ErrNoRows) {
		return -1, password, nil
	}

	return id, password, err
}
