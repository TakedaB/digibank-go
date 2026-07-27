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

func (r *AccountRepository) FindByID(id string) (*model.Account, error) {
	query := `SELECT id, owner, balance FROM accounts WHERE id = $1`

	var account model.Account
	err := r.db.QueryRow(query, id).Scan(&account.ID, &account.Owner, &account.Balance)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) FindAll() ([]model.Account, error) {
	query := `SELECT id, owner, balance FROM accounts ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []model.Account
	for rows.Next() {
		var account model.Account
		if err := rows.Scan(&account.ID, &account.Owner, &account.Balance); err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
