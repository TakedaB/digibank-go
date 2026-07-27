package main

import (
	"log"
	"net/http"

	"github.com/TakedaB/digibank-go/internal/handler"
	"github.com/TakedaB/digibank-go/internal/repository"
)

func main() {
	db := repository.NewPostgresConnection()
	defer db.Close()

	accountRepo := repository.NewAccountRepository(db)
	accountHandler := handler.NewAccountHandler(accountRepo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheckHandler)
	mux.HandleFunc("POST /accounts", accountHandler.Create)
	mux.HandleFunc("GET /accounts/{id}", accountHandler.FindByID)
	mux.HandleFunc("GET /accounts", accountHandler.FindAll)

	log.Println("servidor rodando na porta 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
