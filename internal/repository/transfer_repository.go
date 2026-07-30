package repository

import "database/sql"

type TransferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) *TransferRepository {
	return &TransferRepository{db: db}
}

func (r *TransferRepository) IsProcessed(transferID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM processed_transfers WHERE transfer_id = $1)`

	var exists bool
	err := r.db.QueryRow(query, transferID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TransferRepository) MarkAsProcessed(transferID string) error {
	query := `INSERT INTO processed_transfers (transfer_id) VALUES ($1)`

	_, err := r.db.Exec(query, transferID)
	return err
}

func (r *TransferRepository) CreatePending(id, fromAccountID, toAccountID string, amount int64) error {
	query := `
		INSERT INTO transfers (id, from_account_id, to_account_id, amount, status)
		VALUES ($1, $2, $3, $4, 'pending')
	`

	_, err := r.db.Exec(query, id, fromAccountID, toAccountID, amount)
	return err
}

func (r *TransferRepository) UpdateStatus(id, status string) error {
	query := `UPDATE transfers SET status = $1, updated_at = NOW() WHERE id = $2`

	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *TransferRepository) GetStatus(id string) (string, error) {
	query := `SELECT status FROM transfers WHERE id = $1`

	var status string
	err := r.db.QueryRow(query, id).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}
