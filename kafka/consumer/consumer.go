package consumer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/joho/godotenv/autoload"
	"kafgo/kafka/consumer/variables"
	"kafgo/pkg/logging"
	"sync"
)

type Consumer struct {
	Id            int
	kafkaConsumer *kafka.Consumer
	logger        *logging.Logger
}

func NewConsumer(consumerId int, wg *sync.WaitGroup, logger *logging.Logger) (*Consumer, error) {

	conf := kafka.ConfigMap{
		"bootstrap.servers": variables.KafkaBootstrapServers,
		"group.id":          variables.KafkaGroupId,
		"auto.offset.reset": "earliest",
	}

	consumer, err := kafka.NewConsumer(&conf)
	if err != nil {
		return nil, err
	}

	logger.Printf("New Consumer_%d is starting...", consumerId)

	wg.Done()
	return &Consumer{consumerId, consumer, logger}, nil
}

func (c *Consumer) SubscribeTopic(topic []string) {
	//logger := c.logger.Logger

	err := c.kafkaConsumer.SubscribeTopics(topic, nil)
	if err != nil {
		return
	}
}

func (c *Consumer) ReadMessage(ch chan string) error {
	logger := c.logger.Logger

	for {
		msg, err := c.kafkaConsumer.ReadMessage(-1)
		if err != nil {
			// ToDo залогировать ошибку в бесконечном цикле
			return err
		} else {
			logger.Printf("Consumer_%d Message on %s: %s\n", c.Id, msg.TopicPartition, string(msg.Value))
			ch <- string(msg.Value)
		}
	}

	c.kafkaConsumer.Close()
	logger.Println("Consumer done")

	return nil
}
