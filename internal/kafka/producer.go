package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type TransferProducer struct {
	writer *kafka.Writer
}

func NewTransferProducer(brokerAddress string) *TransferProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokerAddress),
		Topic:    "transfers",
		Balancer: &kafka.LeastBytes{},
	}

	return &TransferProducer{writer: writer}
}

type TransferMessage struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
}

func (p *TransferProducer) Publish(ctx context.Context, msg TransferMessage) error {
	value, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Value: value,
	})
}
