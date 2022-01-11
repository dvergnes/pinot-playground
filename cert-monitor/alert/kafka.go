// Copyright (c) 2022 Denis Vergnes
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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
