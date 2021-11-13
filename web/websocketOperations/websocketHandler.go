package websocketOperations

import (
	"fmt"
	"net/http"

	"github.com/vitalicher97/websocket-api-gateway/external/bitmex"
	"github.com/vitalicher97/websocket-api-gateway/pkg/websocket"
)

func ServeWs(pool *websocket.Pool, bitmexClient *bitmex.WebsocketClient, w http.ResponseWriter, r *http.Request) {
	fmt.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn:         conn,
		Pool:         pool,
		BitmexClient: bitmexClient,
		Subscription: make(map[string]struct{}),
	}

	pool.Register <- client
	client.Read()
}
