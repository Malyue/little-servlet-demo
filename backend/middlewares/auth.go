package middlewares

import (
	"backend/pkg"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

func JWTMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//q := &cmn.ServiceCtx{
		//	R: r,
		//	W: w,
		//
		//	Ep: cmn.Services[r.URL.Path],
		//
		//	Msg: &cmn.ReplyProto{
		//		API: r.URL.Path,
		//		Method: r.Method,
		//	}
		//	BeginTime: time.Now(),
		//}
		//context.WithValue()
		//
		//if r.URL.Path == '/api/login' {
		//	next.ServeHTTP(w, r)
		//}
		//token一般放在Header里的Authorization里，并且以Bearer开头
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			//log.Println("用户未登录")
			zap.L().Error("用户未登录")
			return
		}

		parts := strings.SplitN(authorization, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer ") {
			//log.Println("token非法")
			zap.L().Error("token非法")
			return
		}

		userClaim, err := pkg.ParseToken(parts[1])

		if err != nil {
			//log.Println("token解析失败")
			zap.L().Error("token解析失败", zap.Error(err))
			return
		}

		r.Header.Set("userid", userClaim.Id)
		next.ServeHTTP(w, r)
		//r.Header.Set("userid",)
	})
}
