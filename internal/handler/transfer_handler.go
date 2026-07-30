package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	"github.com/TakedaB/digibank-go/internal/kafka"
	"github.com/TakedaB/digibank-go/internal/repository"
)

type TransferHandler struct {
	producer *kafka.TransferProducer
	repo     *repository.TransferRepository
}

func NewTransferHandler(producer *kafka.TransferProducer, repo *repository.TransferRepository) *TransferHandler {
	return &TransferHandler{producer: producer, repo: repo}
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

	if err := h.repo.CreatePending(transferID, req.FromAccountID, req.ToAccountID, req.Amount); err != nil {
		http.Error(w, "erro ao registrar transferência", http.StatusInternalServerError)
		return
	}

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

func (h *TransferHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	status, err := h.repo.GetStatus(id)
	if err != nil {
		http.Error(w, "transferência não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"transfer_id": id,
		"status":      status,
	})
}
