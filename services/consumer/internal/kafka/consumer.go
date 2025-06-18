package kafka

import (
	"context"
	"distributed-web-scrapper/services/consumer/internal/storage"
	"distributed-web-scrapper/services/consumer/internal/validation"
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

type Consumer struct {
	consumer *kafka.Consumer
	logger   *zap.Logger
}

func NewConsumer(brokers []string, groupID string) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers[0],
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}
	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	logger, _ := zap.NewProduction()
	return &Consumer{consumer: consumer, logger: logger}, nil
}

func (c *Consumer) Consume(ctx context.Context, topics []string, store *storage.PostgresStorage) {
	tracer := otel.Tracer("kafka-consumer")
	if err := c.consumer.SubscribeTopics(topics, nil); err != nil {
		c.logger.Fatal("failed to subscribe to topics", zap.Error(err))
	}

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("shutting down kafka consumer")
			return
		default:
			msg, err := c.consumer.ReadMessage(-1)
			if err != nil {
				c.logger.Error("failed to read message", zap.Error(err))
				continue
			}

			_, span := tracer.Start(ctx, "consume-message")
			var data map[string]interface{}
			if err := json.Unmarshal(msg.Value, &data); err != nil {
				c.logger.Error("failed to unmarshal message", zap.Error(err))
				span.End()
				continue
			}

			if err := validation.ValidateData(data); err != nil {
				c.logger.Error("validation failed", zap.Error(err))
				span.End()
				continue
			}

			if err := store.Save(data); err != nil {
				c.logger.Error("failed to save data", zap.Error(err))
			}
			span.End()
		}
	}
}

func (c *Consumer) Close() {
	c.consumer.Close()
}