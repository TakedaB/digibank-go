package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/TakedaB/digibank-go/internal/repository"
	"github.com/TakedaB/digibank-go/internal/service"
	"github.com/segmentio/kafka-go"
)

type TransferConsumer struct {
	reader          *kafka.Reader
	transferService *service.TransferService
	transferRepo    *repository.TransferRepository
}

func NewTransferConsumer(brokerAddress string, transferService *service.TransferService, transferRepo *repository.TransferRepository) *TransferConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{brokerAddress},
		Topic:   "transfers",
		GroupID: "transfer-processor",
	})

	return &TransferConsumer{
		reader:          reader,
		transferService: transferService,
		transferRepo:    transferRepo,
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

		alreadyProcessed, err := c.transferRepo.IsProcessed(msg.TransferID)
		if err != nil {
			log.Println("erro ao checar idempotência:", err)
			continue
		}

		if alreadyProcessed {
			log.Printf("transferência %s já processada, ignorando (mensagem duplicada)\n", msg.TransferID)
			continue
		}

		log.Printf("processando transferência: %+v\n", msg)

		if err := c.transferService.Transfer(msg.FromAccountID, msg.ToAccountID, msg.Amount); err != nil {
			log.Println("erro ao processar transferência:", err)
			continue
		}

		if err := c.transferRepo.MarkAsProcessed(msg.TransferID); err != nil {
			log.Println("erro ao marcar transferência como processada:", err)
			continue
		}

		log.Println("transferência processada com sucesso")
	}
}
