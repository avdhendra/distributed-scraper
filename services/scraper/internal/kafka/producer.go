package kafka

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

type Producer struct {
	producer *kafka.Producer
	logger   *zap.Logger
}

func NewProducer(brokers []string) (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
	}
	producer, err := kafka.NewProducer(config)
	if err != nil {
		return nil, err
	}

	logger, _ := zap.NewProduction()
	return &Producer{producer: producer, logger: logger}, nil
}

func (p *Producer) Publish(topic string, data interface{}) error {
	msgBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: int32(kafka.PartitionAny)},
		Value:          msgBytes,
	}, nil)
}

func (p *Producer) Close() {
	p.producer.Close()
}