package handler

import (
	"github.com/gin-gonic/gin"
	"kafgo/kafka/consumer"
	"kafgo/kafka/producer"
	"kafgo/pkg/logging"
)

type Handler struct {
	logger      *logging.Logger
	producers   []*producer.Producer
	consumers   map[int]*consumer.Consumer
	messageChan chan string
}

func NewHandler(logger *logging.Logger) *Handler {
	h := &Handler{
		logger:      logger,
		producers:   make([]*producer.Producer, 0),
		consumers:   make(map[int]*consumer.Consumer),
		messageChan: make(chan string, 10000),
	}

	//h.init()
	return h
}

func (h *Handler) init() {
	for i := 1; i < 4; i++ {
		prod, err := producer.NewProducer(int32(i), h.logger)
		if err != nil {
			h.logger.Error(err)
		}
		h.producers = append(h.producers, prod)
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
