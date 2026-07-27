package service

import (
	"errors"

	"github.com/TakedaB/digibank-go/internal/repository"
)

var ErrInsufficientBalance = errors.New("saldo insuficiente")

type TransferService struct {
	accountRepo *repository.AccountRepository
}

func NewTransferService(accountRepo *repository.AccountRepository) *TransferService {
	return &TransferService{accountRepo: accountRepo}
}

func (s *TransferService) Transfer(fromID, toID string, amount int64) error {
	fromAccount, err := s.accountRepo.FindByID(fromID)
	if err != nil {
		return err
	}

	toAccount, err := s.accountRepo.FindByID(toID)
	if err != nil {
		return err
	}

	if fromAccount.Balance < amount {
		return ErrInsufficientBalance
	}

	newFromBalance := fromAccount.Balance - amount
	newToBalance := toAccount.Balance + amount

	if err := s.accountRepo.UpdateBalance(fromID, newFromBalance); err != nil {
		return err
	}

	if err := s.accountRepo.UpdateBalance(toID, newToBalance); err != nil {
		return err
	}
	return nil
}
