package alert

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
)

func NewKafkaNotifier(topic string, producer sarama.SyncProducer) Notifier {
	return &kafkaNotifier{
		topic: topic,
		producer: producer,
	}
}

type kafkaNotifier struct {
	topic    string
	producer sarama.SyncProducer
}

func (k *kafkaNotifier) Send(alert Alert) error {
	data, err := json.Marshal(alert)
	if err != nil {
		return fmt.Errorf("failed to marshal alert in JSON: %w", err)
	}
	_, _, err = k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.ByteEncoder(data),
	})
	if err != nil {
		return fmt.Errorf("failed to deliver alert: %w", err)
	}
	return nil
}

func (k *kafkaNotifier) Close() error {
	if err := k.producer.Close(); err != nil {
		return fmt.Errorf("failed to close producer: %w", err)
	}
	return nil
}
