package dao

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

var dbPool *pgxpool.Pool

func InitDbConn(parseConfig string) {
	zap.L().Info(parseConfig)
	config, err := pgxpool.ParseConfig(parseConfig)
	if err != nil {
		zap.L().Error("parse the pg config error", zap.Error(err))
		return
	}
	dbPool, err = pgxpool.ConnectConfig(context.Background(), config)
	//处理连接失败
	if err != nil {
		zap.L().Error("conn the pgdb error", zap.Error(err))
		return
	}

	fmt.Println("conn the pgdb success")
}

func Connect() (*pgxpool.Conn, error) {
	if dbPool == nil {
		return nil, errors.New("dbPool is nil")
	}
	conn, err := dbPool.Acquire(context.Background())
	return conn, err
}

func Close(conn *pgxpool.Conn) (err error) {
	if conn == nil {
		return errors.New("the conn is nil")
	}
	conn.Release()
	return nil
}
