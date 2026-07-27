package repository

import (
	"database/sql"

	"github.com/TakedaB/digibank-go/internal/model"
)

type AccountRepository struct {
	db *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(owner string, balance int64) (*model.Account, error) {
	query := `
		INSERT INTO accounts (owner, balance)
		VALUES ($1, $2)
		RETURNING id, owner, balance
	`

	var account model.Account
	err := r.db.QueryRow(query, owner, balance).Scan(&account.ID, &account.Owner, &account.Balance)
	if err != nil {
		return nil, err
	}

	return &account, nil
}
