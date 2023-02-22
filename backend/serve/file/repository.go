package file

import (
	"backend/dao"
	"context"
	"errors"
	"github.com/jackc/pgx/v4"
	"go.uber.org/zap"
	"strings"
)

type FileList struct {
	Fid      int    `json:"fid"`
	FileName string `json:"filename"`
}

func InsertFile(fileName string, filePath string) error {
	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return err
	}
	defer dao.Close(conn)
	s := `insert into t_file(filename,filepath) VALUES($1,$2)`

	_, err = conn.Exec(context.Background(), s, fileName, filePath)
	return err
}

func SelectFileByid(fid string) (filename string, filepath string, err error) {
	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return "", "", err
	}
	defer dao.Close(conn)
	s := `select filename,filepath from t_file where fid = $1`
	rc := conn.QueryRow(context.Background(), s, fid)
	err = rc.Scan(&filename, &filepath)

	return
}

func ShowFileList(page int, count int) (fileList []FileList, err error) {
	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return nil, err
	}
	defer dao.Close(conn)

	s := `select fid,filename from t_file limit $1 offset $2`
	rows, err := conn.Query(context.Background(), s, count, (page-1)*count)

	if errors.Is(err, pgx.ErrNoRows) {
		return fileList, nil
	}
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var fid int
	var filename string
	for rows.Next() {
		err = rows.Scan(&fid, &filename)
		if err != nil {
			return
		}
		filename = strings.TrimSpace(filename)
		file := FileList{Fid: fid, FileName: filename}
		fileList = append(fileList, file)
	}

	return
}

func DeleteFile(fid string) (err error) {
	conn, err := dao.Connect()
	if err != nil {
		zap.L().Error("获得数据库连接失败", zap.Error(err))
		return
	}
	defer dao.Close(conn)

	s := `delete from t_file where fid = $1`
	_, err = conn.Exec(context.Background(), s, fid)
	return err
}
