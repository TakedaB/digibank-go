package service_test

import (
	"testing"

	"github.com/TakedaB/digibank-go/internal/repository"
	"github.com/TakedaB/digibank-go/internal/service"
)

func TestTransfer_InsufficientBalance(t *testing.T) {
	db := repository.NewPostgresConnection()
	defer db.Close()

	accountRepo := repository.NewAccountRepository(db)
	transferService := service.NewTransferService(accountRepo)

	//cria duas contas de teste com saldo baixo
	fromAccount, err := accountRepo.Create("Conta Teste Origem", 100)
	if err != nil {
		t.Fatalf("erro ao criar conta de origem: %v", err)
	}
	t.Cleanup(func() { accountRepo.Delete(fromAccount.ID) })

	toAccount, err := accountRepo.Create("Conta Teste Destino", 0)
	if err != nil {
		t.Fatalf("erro ao criar conta de destino: %v", err)
	}
	t.Cleanup(func() { accountRepo.Delete(toAccount.ID) })

	//tenta transferir mais do que o saldo permite
	err = transferService.Transfer(fromAccount.ID, toAccount.ID, 99999)

	if err != service.ErrInsufficientBalance {
		t.Errorf("esperava ErrInsufficientBalance, recebeu: %v", err)
	}
}

func TestTransfer_Success(t *testing.T) {
	db := repository.NewPostgresConnection()
	defer db.Close()

	accountRepo := repository.NewAccountRepository(db)
	transferService := service.NewTransferService(accountRepo)

	fromAccount, err := accountRepo.Create("Conta Teste Origem 2", 10000)
	if err != nil {
		t.Fatalf("erro ao criar conta de origem: %v", err)
	}
	t.Cleanup(func() { accountRepo.Delete(fromAccount.ID) })

	toAccount, err := accountRepo.Create("Conta Teste Destino 2", 0)
	if err != nil {
		t.Fatalf("erro ao criar conta de destino: %v", err)
	}
	t.Cleanup(func() { accountRepo.Delete(toAccount.ID) })

	err = transferService.Transfer(fromAccount.ID, toAccount.ID, 3000)
	if err != nil {
		t.Fatalf("esperava sucesso, recebeu erro: %v", err)
	}

	updatedFrom, err := accountRepo.FindByID(fromAccount.ID)
	if err != nil {
		t.Fatalf("erro ao buscar conta de origem atualizada: %v", err)
	}

	updatedTo, err := accountRepo.FindByID(toAccount.ID)
	if err != nil {
		t.Fatalf("erro ao buscar conta de destino atualizada: %v", err)
	}

	if updatedFrom.Balance != 7000 {
		t.Errorf("esperava saldo de origem 7000, recebeu %d", updatedFrom.Balance)
	}

	if updatedTo.Balance != 3000 {
		t.Errorf("esperava saldo de destino 3000, recebeu %d", updatedTo.Balance)
	}
}
