package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kafgo/kafka/consumer"
	"kafgo/kafka/consumer/variables"
	"kafgo/kafka/producer"
	"net/http"
	"sync"
	"time"
)

var (
	prod *producer.Producer
	cons *consumer.Consumer
)

func (h *Handler) createProducer(ctx *gin.Context) {
	logger := h.logger.Logger

	var err error
	prod, err = producer.NewProducer(h.logger)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"producer": &logger,
	})
}

func (h *Handler) createConsumer(ctx *gin.Context) {
	logger := h.logger.Logger

	consumerId, err := getId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	cons, err = consumer.NewConsumer(consumerId, &wg, h.logger)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	wg.Wait()

	cons.SubscribeTopic([]string{variables.KafkaTopic})

	ctx.JSON(http.StatusOK, gin.H{
		"consumer": &logger,
	})
}

func (h *Handler) postMessage(ctx *gin.Context) {
	logger := h.logger.Logger

	msg := make(chan string)

	go func() {
		err := cons.ReadMessage(msg)
		if err != nil {
			newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
			return
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			prod.SendMessage(fmt.Sprintf("from producer_%d, i: %d", 1, i))
			time.Sleep(4 * time.Second)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			prod.SendMessage(fmt.Sprintf("from producer_%d, i: %d", 2, i))
			time.Sleep(4 * time.Second)
		}
	}()

	for m := range msg {
		logger.Println("Message :", m)
	}
}
