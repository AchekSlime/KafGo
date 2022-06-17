package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
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

	go h.writing(ws)
}

func (h *Handler) writing(conn *websocket.Conn) {
	for {
		select {
		case message := <-h.messageChan:
			conn.WriteMessage(websocket.TextMessage, []byte(message))
		}
	}
}
