package websocket

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"

	clientBitmex "github.com/vitalicher97/websocket-api-gateway/external/bitmex"
	serviceBitmex "github.com/vitalicher97/websocket-api-gateway/service/bitmex"
)

type Client struct {
	ID           string
	Conn         *websocket.Conn
	Pool         *Pool
	BitmexClient *clientBitmex.WebsocketClient
	Subscription map[string]struct{}
}

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

		command := new(serviceBitmex.Command)
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

			_ = serviceBitmex.CommandExecution(c.BitmexClient, command)
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
