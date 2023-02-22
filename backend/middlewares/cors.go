package middlewares

import (
	"log"
	"net/http"
)

func CorsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := r.Method
		origin := r.Header.Get("Origin")

		if origin != "" {
			//接收客户端发送的origin
			w.Header().Set("Access-Control-Allow-Origin", "*")
			//服务器支持所有跨域请求的方法
			w.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIOMS,PUT,DELETE,UPDATE")
			//允许跨域设置可以返回其他字段，跨域自定义字段
			w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Length,X-CSRF-Token,Token,session")
			//允许浏览器可以解析的头部
			w.Header().Set("Access-Control-Expose-Headers", "Content-Lenght,Access-Control-Allow-Origin,Access-Control-Allow-Headers")
			//设置缓存时间
			w.Header().Set("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization") //自定义 Header
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			//c.AbortWithStatus(http.StatusNoContent)
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
