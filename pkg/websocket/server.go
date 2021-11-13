package websocket

import (
	"github.com/vitalicher97/websocket-api-gateway/external/bitmex"
)

func ListenBitmex(pool *Pool, w *bitmex.WebsocketClient) {
	for {
		response, err := w.ReadMessage()
		if err != nil {
			continue
		}

		if _, ok := <-pool.BitmexBroadcast; ok {
			pool.BitmexBroadcast <- response
		}
	}

}
