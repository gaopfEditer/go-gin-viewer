package middleware

import (
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProxyConfig 定义代理配置
type ProxyConfig struct {
	TargetURL   string
	PathPrefix  string
	StripPrefix bool
}

// CreateProxyMiddleware 创建一个代理中间件
func CreateProxyMiddleware(config ProxyConfig) gin.HandlerFunc {
	targetURL, err := url.Parse(config.TargetURL)
	if err != nil {
		panic(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, config.PathPrefix) {
			// 修改请求的Host头
			c.Request.Host = targetURL.Host

			// 如果需要，去掉路径前缀
			if config.StripPrefix {
				c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, config.PathPrefix)
			}

			// 使用反向代理处理请求
			proxy.ServeHTTP(c.Writer, c.Request)
			c.Abort()
		}
	}
}
