package utils

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func (p *Producer) SendMessage(key, value string) error {
	msg := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
		Time:  time.Now(),
	}

	return p.Writer.WriteMessages(nil, msg)
}

func (p *Producer) Close() {
	if err := p.Writer.Close(); err != nil {
		fmt.Println("failed to close write: ", err)
	}
}

func NewProducer(brokerAddress, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &Producer{Writer: writer}
}
