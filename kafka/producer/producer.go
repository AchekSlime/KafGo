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

func NewProducer() *Producer {
	fmt.Println("Producer is starting...")

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": variables.KafkaBootstrapServers})
	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
		//panic(err)
	}

	prd := &Producer{p}

	go prd.deliveryReport()
	return prd
}

func (p *Producer) deliveryReport() {
	for e := range p.kafkaProducer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				log.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				log.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
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
