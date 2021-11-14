package web

import (
	"github.com/gin-gonic/gin"

	"github.com/vitalicher97/websocket-api-gateway/external/bitmex"
	"github.com/vitalicher97/websocket-api-gateway/pkg/websocket"
	"github.com/vitalicher97/websocket-api-gateway/web/websocketOperations"
)

type RouterComponent struct {
	bitmexClient *bitmex.WebsocketClient
}

func NewComponent(w *bitmex.WebsocketClient) *RouterComponent {
	return &RouterComponent{bitmexClient: w}
}

func (rc *RouterComponent) Router(pool *websocket.Pool, r *gin.Engine) {
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ws", func(c *gin.Context) {
		websocketOperations.ServeWs(pool, rc.bitmexClient, c.Writer, c.Request)
	})
}
