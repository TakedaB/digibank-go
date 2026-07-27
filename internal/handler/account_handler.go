package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TakedaB/digibank-go/internal/repository"
)

type AccountHandler struct {
	repo *repository.AccountRepository
}

func NewAccountHandler(repo *repository.AccountRepository) *AccountHandler {
	return &AccountHandler{repo: repo}
}

type createAccountRequest struct {
	Owner   string `json:"owner"`
	Balance int64  `json:"balance"`
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "corpo de requisição inválido", http.StatusBadRequest)
		return
	}

	account, err := h.repo.Create(req.Owner, req.Balance)
	if err != nil {
		http.Error(w, "erro ao criar conta", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	account, err := h.repo.FindByID(id)
	if err != nil {
		http.Error(w, "conta não encontrada", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.repo.FindAll()
	if err != nil {
		http.Error(w, "erro ao buscar contas", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accounts)
}
