package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/joho/godotenv/autoload"
	"kafgo/kafka/consumer/variables"
	"log"
	"sync"
)

type Consumer struct {
	Id            int
	kafkaConsumer *kafka.Consumer
}

func NewConsumer(consumerId int, wg *sync.WaitGroup) (*Consumer, error) {

	conf := kafka.ConfigMap{
		"bootstrap.servers": variables.KafkaBootstrapServers,
		"group.id":          variables.KafkaGroupId,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(&conf)
	if err != nil {
		return nil, err
	}

	log.Printf("New Consumer_%d is starting...", consumerId)

	wg.Done()
	return &Consumer{consumerId, consumer}, nil
}

func (c *Consumer) SubscribeTopic(topic []string) {
	err := c.kafkaConsumer.SubscribeTopics(topic, nil)
	if err != nil {
		return
	}
}

func (c *Consumer) ReadMessage(ch chan string) error {
	for {
		msg, err := c.kafkaConsumer.ReadMessage(-1)
		if err != nil {
			// ToDo залогировать ошибку в бесконечном цикле
			return err
		} else {
			log.Printf("Consumer_%d Message on %s: %s\n", c.Id, msg.TopicPartition, string(msg.Value))
			ch <- string(msg.Value)
		}
	}

	c.kafkaConsumer.Close()
	log.Println("Consumer done")
	return nil
}
