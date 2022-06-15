package producer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kafgo/kafka/producer/variables"
	"log"
)

type Producer struct {
	kafkaProducer *kafka.Producer
}

func NewProducer() (*Producer, error) {
	fmt.Println("Producer is starting...")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": variables.KafkaBootstrapServers})
	if err != nil {
		return nil, err
	}

	prd := &Producer{p}

	go prd.deliveryReport()
	return prd, nil
}

func (p *Producer) deliveryReport() error {
	for e := range p.kafkaProducer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				// ToDo залогировать ошибку в бесконечном цикле
				return ev.TopicPartition.Error
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}

	return nil
}

func (p *Producer) SendMessage(message string) {

	err := p.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &variables.KafkaTopic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, nil)
	if err != nil {
		log.Printf("Send message failed: %s", err)
	}
	//time.Sleep(3 * time.Second)

	// Wait for message deliveries before shutting down
	//p.Flush(15 * 1000)

	log.Println("Successfully send message")
}
