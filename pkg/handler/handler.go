package handler

import (
	"github.com/gin-gonic/gin"
	"kafgo/pkg/logging"
)

type Handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/producer", h.createProducer)
	router.POST("/consumer", h.createConsumer)
	router.POST("/message", h.postMessage)

	return router
}
