package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"

	"github.com/TakedaB/digibank-go/internal/repository"
	"github.com/TakedaB/digibank-go/internal/service"
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
			log.Println("erro ao ler mensagem do kafka:", err)
			continue
		}

		var msg TransferMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Println("erro ao decodificar mensagem:", err)
			continue
		}

		if err := c.ProcessMessage(msg); err != nil {
			log.Println("erro ao processar transferência:", err)
		}
	}
}

// ProcessMessage aplica a lógica de idempotência e processa a transferência.
// Extraída do loop principal para permitir testes sem depender do Kafka.
func (c *TransferConsumer) ProcessMessage(msg TransferMessage) error {
	alreadyProcessed, err := c.transferRepo.IsProcessed(msg.TransferID)
	if err != nil {
		return err
	}

	if alreadyProcessed {
		log.Printf("transferência %s já processada, ignorando (mensagem duplicada)\n", msg.TransferID)
		return nil
	}

	log.Printf("processando transferência: %+v\n", msg)

	if err := c.transferService.Transfer(msg.FromAccountID, msg.ToAccountID, msg.Amount); err != nil {
		c.transferRepo.UpdateStatus(msg.TransferID, "failed")
		return err
	}

	if err := c.transferRepo.MarkAsProcessed(msg.TransferID); err != nil {
		return err
	}

	if err := c.transferRepo.UpdateStatus(msg.TransferID, "completed"); err != nil {
		return err
	}

	log.Println("transferência processada com sucesso")
	return nil
}
