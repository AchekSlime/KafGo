package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kafgo/kafka/consumer"
	"kafgo/kafka/consumer/variables"
	"kafgo/kafka/producer"
	"net/http"
	"sync"
)

func (h *Handler) createProducer(ctx *gin.Context) {
	//logger := h.logger.Logger

	var err error
	prod, err := producer.NewProducer(h.logger)
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

	var wg sync.WaitGroup
	wg.Add(1)
	cons, err := consumer.NewConsumer(consumerId, &wg, h.logger)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	wg.Wait()

	cons.SubscribeTopic([]string{variables.KafkaTopic})
	h.consumers[consumerId] = cons

	//ctx.JSON(http.StatusOK, gin.H{
	//	"consumer": &logger,
	//})
}

func (h *Handler) postMessage(ctx *gin.Context) {
	for i, v := range h.producers {
		v.SendMessage(fmt.Sprintf("from producer_%d", i))
	}
}
