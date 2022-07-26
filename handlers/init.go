package handlers

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Handler struct {
	db *sql.DB
}

func InitDB(db *sql.DB) *Handler {
	return &Handler{db: db}
}
