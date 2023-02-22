package service

//go:generate go run service-enroll-generate.go -a=annotation:(?P<name>.*)-service

import (
	"backend/cmn"
	"backend/dao"
	"backend/middlewares"
	"backend/pkg"
	"backend/serve/chat"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"net/http"
	"os"
	"regexp"
	"runtime"
	"sort"
	"text/template"
	"time"
)

var (
	z      *zap.Logger
	pgConn *pgxpool.Pool
)

var rIsAPI = regexp.MustCompile(`(?i)^/api/(.*)?$`)

var whiteList []string

func init() {
	//初始化配置
}

//func crasher(ctx context.Context) {
//	q := cmn.GetCtxValue(ctx)
//	q.Stop = true
//	r := recover()
//	if r == nil {
//		return
//	}
//
//	reader := bufio.NewReader(strings.NewReader(string(debug.Stack())))
//
//	n := 7
//	var panicStack []string
//	for i := 0; ; i++ {
//		line, _, err := reader.ReadLine()
//		if i == n || i == n+1 {
//			panicStack = append(panicStack, string(line))
//		}
//		if err != nil || i > n+1 {
//			break
//		}
//	}
//
//	templatePanicString := fmt.Sprintf("_CRLF_%s_CRLF_%s",
//		strings.ReplaceAll(strings.Join(panicStack, "_CRLF_"), "\t", ""), r)

//s := strings.ReplaceAll(templatePanicString, "_CRLF_", "\n\t")
//webString := strings.ReplaceAll(templatePanicString, "_CRLF_", ", ")
//q.Err = fmt.Errorf(webString)
//fmt.Println(s)
//fmt.Println(webString)
//返回错误
//q.RespErr()
//}

func proc(reqPath string, w http.ResponseWriter, r *http.Request) {

	q := &cmn.ServiceCtx{

		R: r,
		W: w,

		Ep: cmn.Services[reqPath],

		Msg: &cmn.ReplyProto{
			API:    r.URL.Path,
			Method: r.Method,
		},
		BeginTime: time.Now(),
	}

	//IsWhiteList := false
	//
	//for index, _ := range whiteList {
	//	if reqPath == whiteList[index] {
	//		IsWhiteList = true
	//		break
	//	}
	//}

	////如果非登录或注册接口，则校验token
	//if !IsWhiteList {
	//	var authorization string
	//
	//	//如果是聊天路由，token获取修改
	//
	//	authorization = r.Header.Get("Authorization")
	//	if authorization == "" {
	//		q.Msg.Status = cmn.CodeNeedLogin
	//		q.Msg.Msg = q.Msg.Status.Msg()
	//		q.Resp()
	//		zap.L().Error("用户未登录")
	//		return
	//	}
	//	userClaim, err := pkg.ParseToken(authorization)
	//	if err != nil {
	//		q.Msg.Status = cmn.CodeErrToken
	//		q.Msg.Msg = q.Msg.Status.Msg()
	//		q.Resp()
	//		zap.L().Error("token解析失败", zap.Error(err))
	//		return
	//	}
	//
	//	q.Userid = userClaim.Id
	//}

	ctx := context.WithValue(context.Background(), cmn.QNearKey, q)
	//fmt.Println(q.Userid, "尝试建立连接")

	cmn.Services[reqPath].Fn(ctx)
}

func wsProc(reqPath string, hub *chat.Hub, w http.ResponseWriter, r *http.Request) {
	q := &cmn.ServiceCtx{

		R: r,
		W: w,

		Ep: cmn.Services[reqPath],

		Msg: &cmn.ReplyProto{
			API:    r.URL.Path,
			Method: r.Method,
		},
		BeginTime: time.Now(),
	}
	if reqPath == "/api/testChat" {
		userid := r.URL.Query()["userid"][0]
		q.Userid = userid
	} else {
		//如果非登录或注册接口，则校验token
		var authorization string

		//如果是聊天路由，token获取修改
		//authorization = r.Header.Get("Sec-WebSocket-Protocol")
		authorization = r.URL.Query()["token"][0]
		if authorization == "" {
			q.Msg.Status = cmn.CodeNeedLogin
			q.Msg.Msg = q.Msg.Status.Msg()
			q.Resp()
			zap.L().Error("用户未登录")
			return
		}
		userClaim, err := pkg.ParseToken(authorization)
		if err != nil {
			q.Msg.Status = cmn.CodeErrToken
			q.Msg.Msg = q.Msg.Status.Msg()
			zap.L().Error("token解析失败", zap.Error(err))
			return
		}

		q.Userid = userClaim.Id
	}

	ctx := context.WithValue(context.Background(), cmn.QNearKey, q)

	chat.Chat(ctx, hub)
}

func WebServe() {
	// 注册组件
	Enroll()
	//初始化日志
	cmn.InitLogger()

	//fmt.Println(viper.GetString())
	//初始化pg数据库
	dbConfig := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", viper.GetString("dbms.postgresql.user"),
		viper.GetString("dbms.postgresql.password"),
		viper.GetString("dbms.postgresql.addr"),
		viper.GetInt32("dbms.postgresql.port"),
		viper.GetString("dbms.postgresql.db"))
	dao.InitDbConn(dbConfig)

	whiteList = viper.GetStringSlice("whiteList")

	//初始化redis

	hub := chat.NewHub()
	go hub.Run()

	//注册路由
	router := mux.NewRouter()

	router.Use(middlewares.CorsMiddleWare)

	var rootExists bool
	var pathList []string

	for k := range cmn.Services {
		if k == "/" {
			rootExists = true
			continue
		}

		//add in the pathList
		pathList = append(pathList, k)
	}

	//sort the pathList
	sort.Strings(pathList)

	if rootExists {
		pathList = append(pathList, "/")
	}
	fmt.Println(pathList)
	for _, k := range pathList {
		k := k

		//判断该路径是否是静态文件服务器
		if cmn.Services[k].IsFileServe {
			sys := runtime.GOOS
			docRoot := viper.GetString("docRoot")
			if sys == "windows" {
				docRoot = "./dist/"
			}
			idxFile := docRoot + "/index.html"

			//拦截/static/的请求，将http.Dir中路径替换掉stripPrefix中的路径
			router.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
				t, err := template.ParseFiles(idxFile)
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
				err = t.Execute(w, nil)
				if err != nil {
					fmt.Println(err)
					os.Exit(-1)
				}
			})
			router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir(docRoot+"/assets/"))))
		}

		if cmn.Services[k].IsWebSocket == true {
			router.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
				wsProc(k, hub, w, r)
			})
		}
		router.HandleFunc(k, func(w http.ResponseWriter, r *http.Request) {
			proc(k, w, r)
		})
	}
	port := viper.GetString("port")
	err := http.ListenAndServe(port, router)
	fmt.Println(err)
}
