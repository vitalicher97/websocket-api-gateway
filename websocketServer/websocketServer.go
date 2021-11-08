package websocketServer

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/vitalicher97/websocket-api-gateway/external/bitmex"
)

var upgrader = websocket.Upgrader{}

type (
	ResponseServer struct {
		Timestamp time.Time   `json:"timestamp,omitempty"`
		Symbol    string      `json:"symbol,omitempty"`
		LastPrice json.Number `json:"lastprice,omitempty"`
	}

	Clients struct {
		bitmexClient *bitmex.WebsocketClient
	}
)

func NewClients(w *bitmex.WebsocketClient) *Clients {
	return &Clients{bitmexClient: w}
}

func (c *Clients) RedirectBitmex(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer conn.Close()
	for {
		responseMessage, err := c.bitmexClient.ReadMessage()
		if err != nil {
			return
		}
		for _, data := range responseMessage.Data {
			responseServer := ResponseServer{
				Timestamp: data.Timestamp,
				Symbol:    data.Symbol,
				LastPrice: data.LastPrice,
			}

			if len(responseServer.LastPrice) != 0 {
				err = conn.WriteJSON(responseServer)
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}
	}
}
