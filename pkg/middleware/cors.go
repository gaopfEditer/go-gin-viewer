package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin")
		if origin != "" {

			//局域网内运行的软件可使用*
			//web项目允许127.0.0.1、localhost、域名

			// 允许的来源列表
			allowedOrigins := []string{"http://localhost", "http://127.0.0.1", "https://www.cimagecloud.com"}

			// 检查请求的来源是否在允许列表中
			for _, allowedOrigin := range allowedOrigins {
				if strings.HasPrefix(origin, allowedOrigin) {
					// 使用通配符允许所有端口
					context.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
			//context.Header("Access-Control-Allow-Origin", "*")
			context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			context.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			context.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			context.Header("Access-Control-Expose-Headers", "Authorization")
			context.Header("Access-Control-Allow-Credentials", "true")

		}
		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}
		context.Next()
	}
}
