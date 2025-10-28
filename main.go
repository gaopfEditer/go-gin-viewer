package main

// 使用 _引入依赖项在main函数执行会直接调用init函数
import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"cambridge-hit.com/gin-base/activateserver/app/router"
	_ "cambridge-hit.com/gin-base/activateserver/docs" // swagger
	"cambridge-hit.com/gin-base/activateserver/initializer"
	"cambridge-hit.com/gin-base/activateserver/pkg/middleware"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/ginstatic"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"
	"cambridge-hit.com/gin-base/activateserver/resource"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	initializer.InitAll()

	if resource.Conf.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化引擎
	r := gin.New()
	// 设置中间件
	setMiddleware(r)
	// 设置路由
	setRouter(r)
	// 启动服务(使用goroutine解决服务启动时程序阻塞问题)
	go func() {
		r.Run(fmt.Sprintf("0.0.0.0:%v", resource.Conf.ServerPort))
	}()
	// 监听信号
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-signals:
		// 释放资源
		logger.Sync()
		fmt.Println("[GIN-QuickStart] 程序关闭，释放资源")
		return
	}
}

func setMiddleware(r *gin.Engine) {
	// 日志
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 局域网项目跨域
	r.Use(middleware.Cors())

	// 注册JWT认证中间件
	r.Use(middleware.JwtAuth(resource.Conf.App.ApiPrefix))

	// 注册防抖
	r.Use(middleware.ThrottleMiddleware())
}

func setRouter(r *gin.Engine) {
	// 注册路由
	apiGroup := r.Group(resource.Conf.App.ApiPrefix)
	{
		for _, f := range router.Routers {
			f(apiGroup)
		}
	}

	// Swagger UI 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 设置静态文件夹
	staticDir := "resource/static/web"

	// 检查文件夹是否存在
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		log.Fatalf("Static directory does not exist: %s", staticDir)
	}

	// 使用 Gin 提供静态文件服务
	r.Use(ginstatic.Serve("/", ginstatic.LocalFile(staticDir, true)))
	r.NoRoute(func(context *gin.Context) {
		accept := context.GetHeader("Accept")
		flag := strings.Contains(accept, "text/html")
		if flag {
			content, err := os.ReadFile(staticDir + "/index.html")
			if (err) != nil {
				context.Writer.WriteHeader(404)
				context.Writer.WriteString("Not Found")
				return
			}
			context.Writer.WriteHeader(200)
			context.Writer.Header().Add("Accept", "text/html")
			context.Writer.Write(content)
			context.Writer.Flush()
		}
	})
	fmt.Printf("[GIN-QuickStart] 接口文档地址：http://localhost:%v/swagger/index.html\n", resource.Conf.ServerPort)
	fmt.Printf("[GIN-QuickStart] 前端页面：http://localhost:%v/\n", resource.Conf.ServerPort)
	fmt.Printf("启动时间:%v\n", time.Now().Format(time.DateTime))
}
