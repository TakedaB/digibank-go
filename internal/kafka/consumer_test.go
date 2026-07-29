package kafka_test

import (
	"testing"

	"github.com/TakedaB/digibank-go/internal/kafka"
	"github.com/TakedaB/digibank-go/internal/repository"
	"github.com/TakedaB/digibank-go/internal/service"
	"github.com/google/uuid"
)

func TestProcessMessage_Idempotency(t *testing.T) {
	db := repository.NewPostgresConnection()
	defer db.Close()

	accountRepo := repository.NewAccountRepository(db)
	transferRepo := repository.NewTransferRepository(db)
	transferService := service.NewTransferService(accountRepo)
	consumer := kafka.NewTransferConsumer("localhost:9092", transferService, transferRepo)

	fromAccount, err := accountRepo.Create("Conta Teste Idempotencia Origem", 10000)
	if err != nil {
		t.Fatalf("erro ao criar conta de origem: %v", err)
	}

	toAccount, err := accountRepo.Create("Conta Teste Idempotencia Destino", 0)
	if err != nil {
		t.Fatalf("erro ao criar conta de destino: %v", err)
	}

	msg := kafka.TransferMessage{
		TransferID:    uuid.NewString(),
		FromAccountID: fromAccount.ID,
		ToAccountID:   toAccount.ID,
		Amount:        1000,
	}

	//primeira chamada: deve processar normalmente
	if err := consumer.ProcessMessage(msg); err != nil {
		t.Fatalf("erro ao processar primeira mensagem: %v", err)
	}

	//segunda chamada: mesma mensagem, deve ser ignorada
	if err := consumer.ProcessMessage(msg); err != nil {
		t.Fatalf("erro ao processar mensagem duplicada: %v", err)
	}

	updatedFrom, err := accountRepo.FindByID(fromAccount.ID)
	if err != nil {
		t.Fatalf("erro ao buscar conta atualizada: %v", err)
	}

	//se a idempotência falhasse, o saldo teria saido debitado duas vezes (8000)
	if updatedFrom.Balance != 9000 {
		t.Errorf("esperava saldo 9000(debito unico), recebeu %d - sinal de que a mensagem foi processada mais de uma vez", updatedFrom.Balance)

	}
}
