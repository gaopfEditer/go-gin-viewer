package middleware

import (
	"bytes"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		if !strings.HasPrefix(path, resource.Conf.App.ApiPrefix) {
			context.Next()
			return
		}

		start := time.Now()
		query := context.Request.URL.RawQuery

		// 解析表单数据里面的数据
		if err := context.Request.ParseForm(); err != nil {
			logger.Error("context.Request.ParseForm()", zap.Error(err))
		}
		// 读取表单数据
		form := context.Request.PostForm.Encode()
		bodyStr := ""

		// 判断请求类型是否是json
		if strings.Contains(context.ContentType(), "application/json") {
			defer context.Request.Body.Close()
			body, _ := ioutil.ReadAll(context.Request.Body)
			//注意：重新赋值必须这样否则无法从context重在获取数据
			context.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			bodyStr += string(body)
		}

		context.Next()

		cost := time.Since(start).Microseconds()
		// 接口耗时 小于1000显示单位毫秒 大于1000显示单位秒
		//costTime := map[bool]string{true: fmt.Sprintf("%vms", cost), false: fmt.Sprintf("%vs", float64(cost)/float64(1000))}[cost < 1000]
		costTime := map[bool]string{true: fmt.Sprintf("%vμs", cost), false: fmt.Sprintf("%vms", float64(cost)/float64(1000))}[cost < 1000]

		//costTime := fmt.Sprintf("%vms", cost)

		logger.Info(path,
			zap.Int("status", context.Writer.Status()),
			zap.String("method", context.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Any("form", form),
			zap.Any("json-body", bodyStr),
			zap.String("ip", context.ClientIP()),
			zap.String("user-agent", context.Request.UserAgent()),
			zap.String("errors", context.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.String("cost", costTime),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(context.Request, false)
				if brokenPipe {
					logger.Error(context.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					context.Error(err.(error)) // nolint: errcheck
					context.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				context.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		context.Next()
	}
}
