package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"kafgo"
	"kafgo/kafka/consumer"
	"kafgo/kafka/consumer/variables"
	"kafgo/kafka/producer"
	"net/http"
	"strconv"
	"time"
)

func (h *Handler) createProducer(ctx *gin.Context) {
	//logger := h.logger.Logger
	producerId, err := getId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	prod, err := producer.NewProducer(int32(producerId), h.logger)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	h.producers = append(h.producers, prod)

	//ctx.JSON(http.StatusOK, gin.H{
	//	"producer": &logger,
	//})
}

func (h *Handler) createConsumer(ctx *gin.Context) {
	//logger := h.logger.Logger

	consumerId, err := getId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if _, ok := h.consumers[consumerId]; ok {
		newWarnResponse(ctx, http.StatusBadRequest, "Consumer already exist")
		return
	}

	cons, err := consumer.NewConsumer(consumerId, h.logger)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	cons.SubscribeTopic([]string{variables.KafkaTopic})
	h.consumers[consumerId] = cons

	go cons.ReadMessage(h.messageChan)

	//ctx.JSON(http.StatusOK, gin.H{
	//	"consumer": &logger,
	//})
}

func (h *Handler) postMessage(ctx *gin.Context) {
	for _, v := range h.producers {
		go push(v)
	}
}

func push(prod *producer.Producer) {
	for j := 0; j < 1; j++ {
		message := &kafgo.MessageDto{
			From:      "producer_" + strconv.Itoa(int(prod.Id)),
			To:        "",
			Message:   "test_" + strconv.Itoa(j),
			Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		}
		data, _ := json.Marshal(message)
		prod.SendMessage(data)
	}
}
