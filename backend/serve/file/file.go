package file

import (
	"backend/cmn"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

//author :{"name":"file","email":"1637901557@qq.com"}
//annotation:file-mgmt-service

func init() {
	//HubManage = NewHub()
	//go HubManage.Run()
	//TODO 判断有无文件夹存在
}

func Enroll(author string) {
	zap.L().Info("file.Enroll called")
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
		Fn: uploadFile,

		Path: "/upload",
		Name: "upload",

		Developer: developer,
	})

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: downloadFile,

		Path: "/download",
		Name: "download",

		Developer: developer,
	})

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: deleteFile,

		Path: "/deleteFile",
		Name: "deleteFile",

		Developer: developer,
	})

	cmn.AddService(&cmn.ServeEndPoint{
		Fn: getFileList,

		Path: "/getFile",
		Name: "getFile",

		Developer: developer,
	})

}

func uploadFile(ctx context.Context) {
	fmt.Println("request")
	q := cmn.GetCtxValue(ctx)
	q.R.ParseMultipartForm(32 << 20)
	//獲取上傳文件
	file, handler, err := q.R.FormFile("file")
	if err != nil {
		zap.L().Error("上传文件失败", zap.Error(err))
		return
	}

	defer file.Close()
	//创建上传目录
	//TODO md5Key校验
	os.Mkdir("./uploadFile", os.ModePerm)
	//随机生成名字和后缀
	fileSuffix := path.Ext(handler.Filename)
	randomFileName, _ := uuid.NewUUID()
	//创建上传文件
	f, err := os.Create("./uploadFile/" + randomFileName.String() + fileSuffix)
	defer f.Close()
	io.Copy(f, file)

	//数据库保存名字和服务器名对应
	err = InsertFile(handler.Filename, randomFileName.String()+fileSuffix)
	if err != nil {
		zap.L().Error("文件存入数据库错误", zap.Error(err))
		return
	}

	q.Msg.Status = cmn.CodeSuccess
	q.Msg.Msg = q.Msg.Status.Msg()
	q.Resp()
}

func downloadFile(ctx context.Context) {
	q := cmn.GetCtxValue(ctx)
	//獲取下载文件
	fid := q.R.URL.Query()["fid"][0]
	filename, filepath, err := SelectFileByid(fid)
	if err != nil {
		zap.L().Error("查找数据库文件失败", zap.Error(err))
		return
	}
	//filepath = strings.TrimSpace(filepath)
	file, err := os.Open("./uploadFile/" + filepath)
	defer file.Close()
	//filename = strings.TrimSpace(filename)
	//设置响应头
	q.W.Header().Add("Content-type", "application/octet-stream")
	q.W.Header().Add("content-disposition", "attachment; filename=\""+filename+"\"")
	//将文件写至responseBody
	_, err = io.Copy(q.W, file)
	if err != nil {
		q.W.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(q.W, "Bad request")
		return
	}

}

func getFileList(ctx context.Context) {
	q := cmn.GetCtxValue(ctx)
	//fid := q.R.URL.Query()["fid"][0]
	page, _ := strconv.Atoi(q.R.URL.Query()["page"][0])
	count, _ := strconv.Atoi(q.R.URL.Query()["count"][0])

	fileList, err := ShowFileList(page, count)
	if err != nil {
		zap.L().Error("获得用户列表失败", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = q.Msg.Status.Msg()
		q.Resp()
		return
	}

	q.Msg.Status = cmn.CodeSuccess
	q.Msg.Msg = q.Msg.Status.Msg()
	q.Msg.Data = fileList
	q.Resp()

}

func deleteFile(ctx context.Context) {
	q := cmn.GetCtxValue(ctx)
	fid := q.R.URL.Query()["fid"][0]
	fmt.Println(fid)
	_, filepath, err := SelectFileByid(fid)
	fmt.Println(filepath)
	if err != nil {
		zap.L().Error("查找数据库文件失败", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = q.Msg.Status.Msg()
		q.Resp()
		return
	}

	//删除数据库中文件
	err = DeleteFile(fid)
	if err != nil {
		zap.L().Error("删除数据库文件失败", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = q.Msg.Status.Msg()
		q.Resp()
		return
	}

	//删除本地文件
	err = os.Remove("./uploadFile/" + filepath)
	if err != nil {
		zap.L().Error("删除本地文件失败", zap.Error(err))
		q.Msg.Status = cmn.CodeServerBusy
		q.Msg.Msg = q.Msg.Status.Msg()
		q.Resp()
		return
	}

	q.Msg.Status = cmn.CodeSuccess
	q.Msg.Msg = cmn.CodeSuccess.Msg()
	q.Resp()
}
