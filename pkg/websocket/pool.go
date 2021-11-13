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
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan serviceBitmex.Command),
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
				if err := client.Conn.WriteJSON(message); err != nil {
					fmt.Println(err)
					return
				}
			}
		case bitmexMessage := <-pool.BitmexBroadcast:
			fmt.Println("Sending message to clients in Pool")
			for client, _ := range pool.Clients {
				for _, data := range bitmexMessage.Data {
					if _, ok := client.Subscription[data.Symbol]; ok {
						if err := client.Conn.WriteJSON(bitmexMessage); err != nil {
							fmt.Println(err)
							return
						}
					}
				}
			}
		}
	}
}
