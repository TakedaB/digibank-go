package model

type TransactionStatus string

const (
	StatusPending   TransactionStatus = "pending"
	StatusCompleted TransactionStatus = "completed"
	StatusFailed    TransactionStatus = "failed"
)

type Transaction struct {
	ID            string
	FromAccountID string
	ToAccountID   string
	Amount        int64
	Status        TransactionStatus
}
