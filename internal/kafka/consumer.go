package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/TakedaB/digibank-go/internal/service"
	"github.com/segmentio/kafka-go"
)

type TransferConsumer struct {
	reader          *kafka.Reader
	transferService *service.TransferService
}

func NewTransferConsumer(brokerAddress string, transferService *service.TransferService) *TransferConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   "transfers",
		GroupID: "transfer-processor",
	})

	return &TransferConsumer{
		reader:          reader,
		transferService: transferService,
	}
}

func (c *TransferConsumer) Start(ctx context.Context) {
	log.Println("consumer de transferências iniciado")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("erro ao ler mensagem do kafka", err)
			continue
		}

		var msg TransferMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Println("erro ao decodificar mensagem:", err)
			continue
		}

		log.Printf("processando transferência: %+v\n", msg)

		if err := c.transferService.Transfer(msg.FromAccountID, msg.ToAccountID, msg.Amount); err != nil {
			log.Println("erro ao processar transferência:", err)
			continue
		}

		log.Println("transferência processada com sucesso")
	}
}
