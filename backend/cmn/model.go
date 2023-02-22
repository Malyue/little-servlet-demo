package cmn

import (
	"context"
	"net/http"
	"time"
)

type ReplyProto struct {
	//Status, 0: success, others: fault
	Status RespCode `json:"status"`

	//Msg, Action result describe by literal
	Msg string `json:"msg,omitempty"`

	//Data, operand
	Data interface{} `json:"data,omitempty"`

	// RowCount, just row count
	RowCount int64 `json:"rowCount,omitempty"`

	//API, call target
	API string `json:"API,omitempty"`

	//Method, using http method
	Method string `json:"method,omitempty"`

	//SN, call order
	SN int `json:"SN,omitempty"`
}

// define the service
type ServeEndPoint struct {
	Developer *ModuleAuthor `json:"developer"`

	//Path required,the service url must be unique
	Path string `json:"path,omitempty"`

	//Fn process function
	Fn func(ctx context.Context) `json:"fn,omitempty"`

	//IsFileServe is static html file service
	//true: as the file service
	//false: call fn for service
	IsFileServe bool `json:"is_file_serve,omitempty"`

	IsWebSocket bool `json:"isWebSocket"`

	Name string
}

type ServiceCtx struct {
	Err  error // error occurred during process
	Stop bool  // should run next process

	Attacker  bool // the requester is an attacker
	WhiteList bool // the request path in white list

	Ep *ServeEndPoint

	W http.ResponseWriter
	R *http.Request

	Msg *ReplyProto

	Channel chan []byte

	Userid string

	//用户访问系统所使用的角色
	Role int64

	BeginTime time.Time
}
