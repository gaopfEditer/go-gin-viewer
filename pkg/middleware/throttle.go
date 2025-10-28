package middleware

import (
	"cambridge-hit.com/gin-base/activateserver/resource"
	"strings"
	"sync"

	"cambridge-hit.com/gin-base/activateserver/pkg/util/req-resp/resp"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var (
	limiters = make(map[string]*rate.Limiter)
	mu       sync.Mutex
)

// needThrottle 需要节流的路由及其限制规则
var needThrottle = map[string]struct {
	seconds  int // 秒数
	requests int // 每秒请求次数
	burst    int // 突发大小
}{
	"/dental/base/sendActivationEmail": {20, 1, 1}, // 每60秒最多1次请求, 突发大小为1(即2次)
}

func ThrottleMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiPrefix := "/dental"
		ip := c.ClientIP()
		path := c.Request.URL.Path
		// 前端页面不管
		if !strings.HasPrefix(path, apiPrefix) {
			c.Next()
			return
		}

		if limiterConfig, exists := needThrottle[path]; exists {
			limiter := getLimiter(ip, path, limiterConfig)
			if !limiter.Allow() {
				resp.Error(c, resource.ERR_SERVER_BUSY)
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func getLimiter(ip, path string, config struct {
	seconds  int
	requests int
	burst    int
}) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	key := ip + ":" + path
	limiter, exists := limiters[key]
	if !exists {
		// 计算速率：每秒允许的请求次数
		rateLimit := rate.Limit(float64(config.requests) / float64(config.seconds))
		limiter = rate.NewLimiter(rateLimit, config.burst)
		limiters[key] = limiter
	}

	return limiter
}
