package service

import (
	"encoding/json"
	"sync"
	"time"

	"cambridge-hit.com/gin-base/activateserver/pkg/util/logger"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	// 允许客户端写入的最大消息大小
	maxMessageSize = 512 * 1024 // 512KB

	// 客户端写入超时时间
	writeWait = 10 * time.Second

	// 客户端读取超时时间
	pongWait = 60 * time.Second

	// 发送 ping 的间隔，必须小于 pongWait
	pingPeriod = (pongWait * 9) / 10

	// 允许客户端关闭连接的时间
	closeGracePeriod = 10 * time.Second
)

// Client 表示一个 WebSocket 客户端连接
type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID int // 用户ID，可以从 JWT token 中获取
}

// Hub 维护所有活跃的客户端连接并广播消息
type Hub struct {
	// 已注册的客户端
	Clients map[*Client]bool

	// 从客户端接收到的消息
	Broadcast chan []byte

	// 注册新客户端
	Register chan *Client

	// 注销客户端
	Unregister chan *Client

	// 互斥锁，保护 Clients map
	mu sync.RWMutex
}

// NewHub 创建新的 Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run 启动 Hub，处理客户端注册、注销和消息广播
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()
			logger.Info("WebSocket客户端已连接", zap.Int("userID", client.UserID), zap.Int("连接数", len(h.Clients)))

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
			}
			h.mu.Unlock()
			logger.Info("WebSocket客户端已断开", zap.Int("userID", client.UserID), zap.Int("连接数", len(h.Clients)))

		case message := <-h.Broadcast:
			h.mu.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					// 如果发送失败，关闭连接
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// BroadcastMessage 广播消息给所有客户端
func (h *Hub) BroadcastMessage(message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	h.Broadcast <- data
	return nil
}

// SendToUser 发送消息给特定用户
func (h *Hub) SendToUser(userID int, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.Clients {
		if client.UserID == userID {
			select {
			case client.Send <- data:
			default:
				// 发送失败，关闭连接
				close(client.Send)
				delete(h.Clients, client)
			}
		}
	}
	return nil
}

// GetClientCount 获取当前连接的客户端数量
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Clients)
}

// ReadPump 从 WebSocket 连接读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket读取错误", zap.Error(err))
			}
			break
		}

		// 处理接收到的消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err == nil {
			logger.Info("收到WebSocket消息", zap.Any("message", msg), zap.Int("userID", c.UserID))
			// 这里可以添加消息处理逻辑
			// 例如：根据消息类型执行不同的操作
		}
	}
}

// WritePump 向 WebSocket 连接写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将队列中的其他消息也发送出去
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WebSocketService WebSocket 服务
type WebSocketService struct {
	Hub *Hub
}

// NewWebSocketService 创建新的 WebSocket 服务
func NewWebSocketService() *WebSocketService {
	hub := NewHub()
	go hub.Run()
	return &WebSocketService{
		Hub: hub,
	}
}

// GetHub 获取 Hub 实例
func (s *WebSocketService) GetHub() *Hub {
	return s.Hub
}
