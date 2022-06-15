package main

import (
	"context"
	"kafgo"
	"kafgo/pkg/handler"
	"kafgo/pkg/logging"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logger := logging.GetLogger()
	handlers := handler.NewHandler(logger)

	/*var wg sync.WaitGroup
	wg.Add(1)
	consumer1 := consumer.NewConsumer(1, &wg)
	//consumer2 := consumer.NewConsumer(2, &wg)
	wg.Wait()

	consumer1.SubscribeTopic([]string{variables.KafkaTopic})
	//consumer2.SubscribeTopic([]string{variables.KafkaTopic})

	msg := make(chan string)

	go consumer1.ReadMessage(msg)
	//go consumer2.ReadMessage(msg)

	producer1 := producer.NewProducer()
	producer2 := producer.NewProducer()

	go func() {
		for i := 0; i < 100; i++ {
			producer1.SendMessage(fmt.Sprintf("from producer_%d, i: %d", 1, i))
			time.Sleep(4 * time.Second)
		}
	}()

	go func() {
		for i := 0; i < 100; i++ {
			producer2.SendMessage(fmt.Sprintf("from producer_%d, i: %d", 2, i))
			time.Sleep(4 * time.Second)
		}
	}()

	for m := range msg {
		log.Println("Message :", m)
	}*/

	srv := new(kafgo.Server)
	go func() {
		if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
			logger.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logger.Info("server Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Errorf("error occured on server shutting down: %s", err.Error())
	}

	logger.Info("server Shutting Down")
}
