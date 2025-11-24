package controller

import (
	"net/http"

	"cambridge-hit.com/gin-base/activateserver/app/service"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/auth"
	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源，生产环境应该检查 Origin
		return true
	},
}

type WebSocketController struct {
	wsService *service.WebSocketService
}

// NewWebSocketController 创建新的 WebSocket 控制器
func NewWebSocketController() *WebSocketController {
	return &WebSocketController{
		wsService: service.NewWebSocketService(),
	}
}

// HandleWebSocket 处理 WebSocket 连接
// @Tags     WebSocket
// @Summary  WebSocket 连接
// @Produce  application/json
// @Router   /activate/ws [get]
func (wc *WebSocketController) HandleWebSocket(c *gin.Context) {
	// 从 JWT token 中获取用户信息（可选）
	userInfo, exists := c.Get("userInfo")
	var userID int
	if exists {
		if uai, ok := userInfo.(auth.UserAuthInfo); ok {
			userID = uai.UserID
		}
	}

	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("WebSocket升级失败", zap.Error(err))
		return
	}

	// 创建客户端
	client := &service.Client{
		Hub:    wc.wsService.GetHub(),
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	// 注册客户端
	client.Hub.Register <- client

	// 启动客户端的读写协程
	go client.WritePump()
	go client.ReadPump()
}

// BroadcastMessage 广播消息给所有连接的客户端
// @Tags     WebSocket
// @Summary  广播消息
// @Produce  application/json
// @Param    message  body      map[string]interface{}  true  "消息内容"
// @Router   /activate/ws/broadcast [post]
func (wc *WebSocketController) BroadcastMessage(c *gin.Context) {
	var message map[string]interface{}
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的消息格式"})
		return
	}

	if err := wc.wsService.GetHub().BroadcastMessage(message); err != nil {
		logger.Error("广播消息失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "广播消息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "消息已广播"})
}

// GetClientCount 获取当前连接的客户端数量
// @Tags     WebSocket
// @Summary  获取连接数
// @Produce  application/json
// @Router   /activate/ws/count [get]
func (wc *WebSocketController) GetClientCount(c *gin.Context) {
	count := wc.wsService.GetHub().GetClientCount()
	c.JSON(http.StatusOK, gin.H{"count": count})
}
