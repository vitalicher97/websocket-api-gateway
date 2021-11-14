package main

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/vitalicher97/websocket-api-gateway/external/bitmex"
	"github.com/vitalicher97/websocket-api-gateway/pkg/websocket"
	"github.com/vitalicher97/websocket-api-gateway/web"
)

func main() {
	urlToBitmexWebsocket := url.URL{
		Scheme: "wss",
		Host:   "ws.testnet.bitmex.com",
		Path:   "/realtime",
	}

	websocketClient := bitmex.NewWebsocketClient(nil, urlToBitmexWebsocket)
	websocketClient, err := websocketClient.SetConnection()
	if err != nil {
		log.Panicln("error websocket connection was not set")
	}

	pool := websocket.NewPool()
	go pool.Start()

	go websocket.ListenBitmex(pool, websocketClient)

	component := web.NewComponent(websocketClient)
	r := gin.Default()
	component.Router(pool, r)
	err = r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	if err != nil {
		log.Panicln("error server could not start: %s", err)
	}

}
