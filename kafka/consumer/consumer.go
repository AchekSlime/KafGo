package consumer

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/joho/godotenv/autoload"
	"kafgo"
	"kafgo/kafka/consumer/variables"
	"kafgo/pkg/logging"
	"strconv"
	"time"
)

type Consumer struct {
	Id            int
	kafkaConsumer *kafka.Consumer
	logger        *logging.Logger
}

func NewConsumer(consumerId int, logger *logging.Logger) (*Consumer, error) {

	conf := kafka.ConfigMap{
		"bootstrap.servers": variables.KafkaBootstrapServers,
		"group.id":          variables.KafkaGroupId + strconv.Itoa(consumerId),
		"auto.offset.reset": "latest",
	}

	consumer, err := kafka.NewConsumer(&conf)
	if err != nil {
		return nil, err
	}

	logger.Printf("NEW Consumer_%d is starting...", consumerId)

	return &Consumer{consumerId, consumer, logger}, nil
}

func (c *Consumer) SubscribeTopic(topic []string) {
	//logger := c.logger.Logger

	err := c.kafkaConsumer.SubscribeTopics(topic, nil)
	if err != nil {
		return
	}
}

func (c *Consumer) ReadMessage(ch chan string) {
	logger := c.logger.Logger

	logger.Printf("Cnosumer_%d: start reading from topic", c.Id)
	for {
		msg, err := c.kafkaConsumer.ReadMessage(-1)
		if err != nil {
			logger.Error(err)
		} else {
			logger.Printf("Consumer_%d: messege: %s on topic: %v\n", c.Id, string(msg.Value), msg.TopicPartition)
			var message kafgo.MessageDto
			json.Unmarshal(msg.Value, &message)
			message.To = "consumer_" + strconv.Itoa(c.Id)
			data, _ := json.Marshal(message)
			ch <- string(data)
		}
		time.Sleep(time.Second)
	}

	//c.kafkaConsumer.Close()
	//logger.Println("Consumer done")
	//
}
