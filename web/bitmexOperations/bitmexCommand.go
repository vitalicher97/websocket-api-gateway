package bitmexOperations

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vitalicher97/websocket-api-gateway/external/bitmex"
	bitmex2 "github.com/vitalicher97/websocket-api-gateway/service/bitmex"
)

type Handler struct {
	bitmexClient *bitmex.WebsocketClient
}

var (
	// ErrInvalidRequest if request is invalid
	ErrInvalidRequest = errors.New("error invalid request")
)

func NewHandler(w *bitmex.WebsocketClient) *Handler {
	return &Handler{bitmexClient: w}
}

func (h *Handler) CommandExecution(c *gin.Context) {
	command := new(bitmex2.Command)
	newDecoder := json.NewDecoder(c.Request.Body)
	newDecoder.DisallowUnknownFields()
	err := newDecoder.Decode(command)
	if err != nil {
		log.Println(ErrInvalidRequest, ":", err)
		errorResponse := bitmex2.ErrorResponse{
			Code:    http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	log.Println(command) // Should be debug level log

	err = bitmex2.CommandExecution(h.bitmexClient, command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"ok": "ok",
	})
	return
}
