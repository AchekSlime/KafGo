package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"kafgo/kafka/consumer"
	"log"
	"net/http"
	"strconv"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) connection(ctx *gin.Context) {
	// ws соединение
	ws, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	_, message, err := ws.ReadMessage()
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	id, err := strconv.Atoi(string(message))
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	cons, ok := h.consumers[id]
	if !ok {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	go read(cons, ws)
}

func read(cons *consumer.Consumer, conn *websocket.Conn) {
	msgs := make(chan string)
	go func() {
		err := cons.ReadMessage(msgs)
		if err != nil {
			log.Println("reading message error: ", err)
			return
		}
	}()

	for {
		select {
		case message := <-msgs:
			conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}
