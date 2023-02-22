package file

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"
)

var index = 0

func UploadFile() {
	index++
	beginTime := time.Now()
	url := "http://localhost:9090/api/upload"
	paramName := "file"
	filePath := "../../dist/index.html"
	//打开要上传的文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("文件上传错误", err)
		return
	}
	defer file.Close()
	body := &bytes.Buffer{}
	//创建一个multipart类型
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filePath)
	if err != nil {
		fmt.Println("post err =", err)
		return
	}
	_, err = io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", url, body)
	//request.Header.Add()
	request.Header.Set("Content-Type", writer.FormDataContentType())
	clt := http.Client{
		Timeout: 5 * time.Second,
	}
	defer clt.CloseIdleConnections()
	res, err := clt.Do(request)
	defer func() {
		res.Body.Close()
		fmt.Println("响应时间", time.Now().Sub(beginTime))
		if time.Now().Sub(beginTime) > 3*time.Second {
			fmt.Println("共有", index, "进程")
			os.Exit(1)
		}
	}()
}

func BenchmarkFile(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			//上传文件
			go UploadFile()
		}
	})
}
