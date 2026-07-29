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
