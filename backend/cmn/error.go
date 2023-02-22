package cmn

type RespCode int

const (
	CodeSuccess    RespCode = 0
	CodeServerBusy RespCode = 1000 + iota
	CodeNeedLogin
	CodeErrLogin
	CodeErrToken
	CodeInvalidParam
	CodeHadUser
)

var codeMsgMap = map[RespCode]string{
	CodeSuccess:      "success",
	CodeServerBusy:   "服务繁忙",
	CodeNeedLogin:    "需要登录",
	CodeErrLogin:     "账号密码错误",
	CodeErrToken:     "非法token",
	CodeInvalidParam: "非法参数",
	CodeHadUser:      "已存在该账号",
}

func (c RespCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		return codeMsgMap[CodeServerBusy]
	}
	return msg
}
