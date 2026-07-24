package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresConnection() *sql.DB {
	connStr := "postgres://postgres:senha123@localhost:5432/digibank?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("erro ao abrir conexão com o banco:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("erro ao conectar com o banco:", err)
	}

	fmt.Println("conectado aoo postgres com sucesso")
	return db

}
