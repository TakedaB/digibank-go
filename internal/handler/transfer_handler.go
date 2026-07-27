package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TakedaB/digibank-go/internal/service"
)

type TransferHandler struct {
	service *service.TransferService
}

func NewTransferHandler(service *service.TransferService) *TransferHandler {
	return &TransferHandler{service: service}
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

	err := h.service.Transfer(req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		if errors.Is(err, service.ErrInsufficientBalance) {
			http.Error(w, "saldo insuficiente", http.StatusUnprocessableEntity)
			return
		}
		http.Error(w, "erro ao processar transferência", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
