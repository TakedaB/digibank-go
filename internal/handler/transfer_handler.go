package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/TakedaB/digibank-go/internal/kafka"
)

type TransferHandler struct {
	producer *kafka.TransferProducer
}

func NewTransferHandler(producer *kafka.TransferProducer) *TransferHandler {
	return &TransferHandler{producer: producer}
}

type transferRequest struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
}

func (h *TransferHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	var req transferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	transferID := uuid.NewString()

	msg := kafka.TransferMessage{
		TransferID:    transferID,
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}

	if err := h.producer.Publish(r.Context(), msg); err != nil {
		http.Error(w, "erro ao processar transferência", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"status":      "processando",
		"transfer_id": transferID,
	})
}
