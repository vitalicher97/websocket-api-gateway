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

		fmt.Printf("Message Received: %+v\n", command)

		if command.Action == "subscribe" {
			if len(command.Symbols) == 0 {
				c.Subscription["ALL"] = struct{}{}
			} else {
				for _, symbol := range command.Symbols {
					c.Subscription[symbol] = struct{}{}
				}
			}

			_ = bitmex2.CommandExecution(c.BitmexClient, command)
		}

		if command.Action == "unsubscribe" {
			if len(command.Symbols) == 0 {
				for key := range c.Subscription {
					delete(c.Subscription, key)
				}
			} else {
				for _, symbol := range command.Symbols {
					delete(c.Subscription, symbol)
				}
			}
		}
	}
}
