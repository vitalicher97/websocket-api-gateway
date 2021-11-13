package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	"github.com/vitalicher97/websocket-api-gateway/external/bitmex"
	bitmex2 "github.com/vitalicher97/websocket-api-gateway/service/bitmex"
)

type Client struct {
	ID           string
	Conn         *websocket.Conn
	Pool         *Pool
	BitmexClient *bitmex.WebsocketClient
	Subscription map[string]struct{}
}

/*type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}*/

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		command := new(bitmex2.Command)
		err = json.Unmarshal(p, command)
		if err != nil {
			log.Println("Invalid Message")
		}

		for _, symbol := range command.Symbols {
			c.Subscription[symbol] = struct{}{}
		}

		fmt.Printf("Message Received: %+v\n", command)
		err = bitmex2.CommandExecution(c.BitmexClient, command)
		if err != nil {
			continue
		}

	}
}
