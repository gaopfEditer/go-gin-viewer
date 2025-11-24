package router

import (
	"cambridge-hit.com/gin-base/activateserver/app/controller"
	"github.com/gin-gonic/gin"
)

func init() {
	Routers = append(Routers, WebSocketRouterRegister)
}

// WebSocketRouterRegister 注册 WebSocket 路由
func WebSocketRouterRegister(r *gin.RouterGroup) {
	wsGroup := r.Group("ws")
	wsController := controller.NewWebSocketController()
	{
		// WebSocket 连接端点
		wsGroup.GET("", wsController.HandleWebSocket)
		// 广播消息端点（需要认证）
		wsGroup.POST("/broadcast", wsController.BroadcastMessage)
		// 获取连接数端点（需要认证）
		wsGroup.GET("/count", wsController.GetClientCount)
	}
}
