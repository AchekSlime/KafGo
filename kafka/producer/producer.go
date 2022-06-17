package producer

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kafgo/kafka/producer/variables"
	"kafgo/pkg/logging"
)

type Producer struct {
	kafkaProducer *kafka.Producer
	logger        *logging.Logger
	Id            int32
}

func NewProducer(producerId int32, logger *logging.Logger) (*Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": variables.KafkaBootstrapServers})
	if err != nil {
		return nil, err
	}

	prd := &Producer{p, logger, producerId}

	go prd.deliveryReport()
	logger.Printf("NEW Producer_%d is started...", prd.Id)
	return prd, nil
}

func (p *Producer) deliveryReport() error {
	logger := p.logger.Logger

	for e := range p.kafkaProducer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				logger.Printf("Delivery err: %v\n", ev.TopicPartition.Error)
				return ev.TopicPartition.Error
			} else {
				logger.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}

	return nil
}

func (p *Producer) SendMessage(data []byte) {
	logger := p.logger.Logger

	err := p.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &variables.KafkaTopic, Partition: kafka.PartitionAny},
		Value:          data,
	}, nil)
	if err != nil {
		logger.Warning("Send message failed: %s", err)
	} else {
		logger.Info("Successfully send message")
	}
	//time.Sleep(3 * time.Second)

	// Wait for message deliveries before shutting down
	//p.Flush(15 * 1000)

}
