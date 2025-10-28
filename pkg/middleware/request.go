package middleware

import (
	"bytes"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"io"

	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func JsonDataMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		data, err := context.GetRawData()
		if err != nil {
			logger.Error("Error reading raw data:", zap.Error(err))
			resp.Error(context, resource.ERR_NO_PERMISSION)
			context.Abort()
			return
		}
		// 恢复请求体流数据
		context.Request.Body = io.NopCloser(bytes.NewBuffer(data))

		// 将解析后的JSON数据存储到Context中
		context.Set("rawRequestData", data)

		context.Next()

	}
}
