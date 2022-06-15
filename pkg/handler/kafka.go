package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kafgo/kafka/consumer"
	"kafgo/kafka/consumer/variables"
	"kafgo/kafka/producer"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	prod *producer.Producer
	cons *consumer.Consumer
)

func (h *Handler) createProducer(ctx *gin.Context) {
	prod = producer.NewProducer()

	ctx.JSON(http.StatusOK, gin.H{
		"producer": &prod,
	})
}

func (h *Handler) createConsumer(ctx *gin.Context) {
	consumerId, err := getId(ctx)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(1)
	cons = consumer.NewConsumer(consumerId, &wg)
	wg.Wait()

	cons.SubscribeTopic([]string{variables.KafkaTopic})

	ctx.JSON(http.StatusOK, gin.H{
		"consumer": &cons,
	})
}

func (h *Handler) postMessage(ctx *gin.Context) {
	msg := make(chan string)

	go cons.ReadMessage(msg)

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
		log.Println("Message :", m)
	}
}
