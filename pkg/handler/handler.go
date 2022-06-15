package handler

import (
	"github.com/gin-gonic/gin"
	"kafgo/kafka/consumer"
	"kafgo/kafka/producer"
	"kafgo/pkg/logging"
)

type Handler struct {
	logger    *logging.Logger
	producers []*producer.Producer
	consumers map[int]*consumer.Consumer
}

func NewHandler(logger *logging.Logger) *Handler {
	return &Handler{
		logger:    logger,
		producers: make([]*producer.Producer, 0),
		consumers: make(map[int]*consumer.Consumer),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/connection", h.connection)
	router.POST("/producer", h.createProducer)
	router.POST("/consumer", h.createConsumer)
	router.POST("/message", h.postMessage)

	return router
}
