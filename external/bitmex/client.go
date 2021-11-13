package bitmex

import (
	"encoding/json"
	"errors"
	"log"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/go-playground/validator.v9"
)

var (
	// ErrSetConnection if failed to set connection
	ErrSetConnection = errors.New("error on set connection with Bitmex service")
	// ErrCommandInvalid if command is invalid
	ErrCommandInvalid = errors.New("error provided command is invalid for Bitmex")
	// ErrSendMessage if message was not sent
	ErrSendMessage = errors.New("error on send message to Bitmex")
	// ErrReadMessage if message can not be read
	ErrReadMessage = errors.New("error on read message from Bitmex")
	// ErrInvalidResponse if response is invalid
	ErrInvalidResponse = errors.New("error response is invalid from Bitmex")
)

type (
	// WebsocketClient is for websocket connection
	WebsocketClient struct {
		wsConn       *websocket.Conn
		urlToConnect string
	}

	// Command is for sending commands to websocket server
	Command struct {
		Op   string   `json:"op" validate:"required"`
		Args []string `json:"args,omitempty"`
	}

	// ResponseMessage is for receiving messages from websocket server
	ResponseMessage struct {
		Data      []Data `json:"data,omitempty"`
	}

	// Data is for nested values
	Data struct {
		Timestamp time.Time   `json:"timestamp,omitempty"`
		Symbol    string      `json:"symbol,omitempty"`
		LastPrice json.Number `json:"lastprice,omitempty"`
	}
)

// NewWebsocketClient to initialize WebsocketClient
func NewWebsocketClient(wsConn *websocket.Conn, urlToConnect url.URL) *WebsocketClient {
	return &WebsocketClient{
		wsConn:       wsConn,
		urlToConnect: urlToConnect.String(),
	}
}

// SetConnection to create a connection with websocket server
func (w *WebsocketClient) SetConnection() (*WebsocketClient, error) {
	if w.wsConn != nil {
		return w, nil
	}

	conn, resp, err := websocket.DefaultDialer.Dial(w.urlToConnect, nil)
	if err != nil {
		log.Println(ErrSetConnection, ":", err)
		return nil, ErrSetConnection
	}

	w.wsConn = conn

	dumpResp, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Println("error on dump of response")
	}

	log.Println(string(dumpResp)) // Should be debug level log

	return w, nil
}

// SendCommand to send command to the websocket server
func (w *WebsocketClient) SendCommand(message Command) error {
	validate := validator.New()
	err := validate.Struct(message)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			log.Println(e)
		}
		return ErrCommandInvalid
	}

	err = w.wsConn.WriteJSON(message)
	if err != nil {
		log.Println("error on write JSON to websocket external server:", err)
		return ErrSendMessage
	}

	return nil
}

// ReadMessage to receive message from websocket server
func (w *WebsocketClient) ReadMessage() (*ResponseMessage, error) {
	msgType, msg, err := w.wsConn.ReadMessage()
	if err != nil {
		log.Println(ErrReadMessage, ":", err)
		return nil, ErrReadMessage
	}

	log.Printf("Type: %s, Message: %s\n", msgType, msg) // Should be debug level log

	responseMessage := new(ResponseMessage)
	err = json.Unmarshal(msg, responseMessage)
	if err != nil {
		log.Println(ErrInvalidResponse, ":", err)
		return nil, ErrInvalidResponse
	}

	return responseMessage, nil
}
