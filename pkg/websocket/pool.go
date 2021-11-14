package websocket

import (
	"fmt"
	"log"

	"github.com/vitalicher97/websocket-api-gateway/external/bitmex"
	serviceBitmex "github.com/vitalicher97/websocket-api-gateway/service/bitmex"
)

type Pool struct {
	Register        chan *Client
	Unregister      chan *Client
	Clients         map[*Client]bool
	Broadcast       chan serviceBitmex.Command
	BitmexBroadcast chan *bitmex.ResponseMessage
}

func NewPool() *Pool {
	return &Pool{
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Clients:         make(map[*Client]bool),
		Broadcast:       make(chan serviceBitmex.Command),
		BitmexBroadcast: make(chan *bitmex.ResponseMessage),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				log.Println("Connected: ", client)
			}
		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			fmt.Println("Size of Connection Pool: ", len(pool.Clients))
			for client, _ := range pool.Clients {
				log.Println("Disconnected: ", client)
			}
		case message := <-pool.Broadcast:
			fmt.Println("Sending message to all clients in Pool")
			for client, _ := range pool.Clients {
				log.Println(client, message)
			}
		case bitmexMessage := <-pool.BitmexBroadcast:
			fmt.Println("Sending message to clients in Pool")
			for client, _ := range pool.Clients {
				for _, data := range bitmexMessage.Data {
					_, all := client.Subscription["ALL"]
					if _, ok := client.Subscription[data.Symbol]; (ok || all) && len(data.LastPrice) != 0 {
						if err := client.Conn.WriteJSON(data); err != nil {
							fmt.Println(err)
							return
						}
					}
				}
			}
		}
	}
}
