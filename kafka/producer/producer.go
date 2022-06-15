package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kafgo/kafka/producer/variables"
	"kafgo/pkg/logging"
)

type Producer struct {
	kafkaProducer *kafka.Producer
	logger        *logging.Logger
}

func NewProducer(logger *logging.Logger) (*Producer, error) {
	logger.Println("Producer is starting...")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": variables.KafkaBootstrapServers})
	if err != nil {
		return nil, err
	}

	prd := &Producer{p, logger}

	go prd.deliveryReport()
	return prd, nil
}

func (p *Producer) deliveryReport() error {
	logger := p.logger.Logger

	for e := range p.kafkaProducer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				// ToDo залогировать ошибку в бесконечном цикле
				return ev.TopicPartition.Error
			} else {
				logger.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}

	return nil
}

func (p *Producer) SendMessage(message string) {
	logger := p.logger.Logger

	err := p.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &variables.KafkaTopic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		logger.Printf("Send message failed: %s", err)
	}
	//time.Sleep(3 * time.Second)

	// Wait for message deliveries before shutting down
	//p.Flush(15 * 1000)

	logger.Println("Successfully send message")
}
