package cmn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"sync"
)

type ModuleAuthor struct {
	Name  string `json:"name"`
	Tel   string `json:"tel"`
	Email string `json:"email"`
	Addi  string `json:"addi"`
}

var (
	Services = make(map[string]*ServeEndPoint)

	serviceMutex sync.Mutex
)

type ctxKey string

const QNearKey = ctxKey("ServiceCtx")

var rIsAPI = regexp.MustCompile(`(?i)^/api/(.*)?$`)

func GetCtxValue(ctx context.Context) (q *ServiceCtx) {
	var err error
	f := ctx.Value(QNearKey)
	if f == nil {
		err = fmt.Errorf(`get nil from ctx.Value["%s]`, fmt.Sprint(QNearKey))
		panic(err.Error())
	}

	var ok bool
	q, ok = f.(*ServiceCtx)
	if !ok {
		err := fmt.Errorf("failer to type assertion for *ServiceCtx")
		panic(err.Error())
	}
	if q == nil {
		err := fmt.Errorf(`ctx.Value["%s"] should be non nil *ServiceCtx`, fmt.Sprint(QNearKey))
		panic(err.Error())
	}

	return
}

func AddService(ep *ServeEndPoint) (err error) {
	for {
		if ep == nil {
			err = errors.New("ep is nil")
			break
		}

		if ep.Path == "" {
			err = errors.New("ep.path empty")
			break
		}

		if ep.IsFileServe {
			//if ep.DocRoot == "" {
			//	err = errors.New("must specify docRoot when ep.isFileServe equal true")
			//	break
			//}

			//if ep.Fn == nil {
			//	ep.Fn = WebFS
			//}
		} else {
			if ep.Fn == nil {
				err = errors.New("must specify fn when ep.isFileServe equal false")
				break
			}

			if !rIsAPI.MatchString(ep.Path) {
				ep.Path = strings.ReplaceAll("/api/"+ep.Path, "//", "/")
			}
		}
		_, ok := Services[ep.Path]
		if ok {
			err = errors.New(fmt.Sprintf("%s[%s] already exists", ep.Path, ep.Name))
		}
		break

	}

	if err != nil {
		//z.Error(err.Error())
		fmt.Errorf(err.Error())
		return
	}

	serviceMutex.Lock()
	defer serviceMutex.Unlock()

	Services[ep.Path] = ep
	return
}

func (v *ServiceCtx) Resp() {
	if v.W == nil || v.Msg == nil {
		zap.L().Error("call respErr with invalid w/msg")
		return
	}

	//将msg结构体转为json字符串
	buf, err := json.Marshal(v.Msg)
	if err != nil {
		zap.L().Error("转换msg错误", zap.Error(err))
		_, _ = fmt.Fprintf(v.W, err.Error())
		return
	}

	v.W.Write(buf)
}

func Resp() {

}
