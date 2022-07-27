package main

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"go-todos-api/handlers"
)

func main() {
	// Database credentials
	const (
		host = "localhost"
		port = 5432
		user = "postgres"
		password = "root"
		dbname = "postgres"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Validate database credentials
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	defer db.Close()

	// Connect database
	if err := db.Ping(); err != nil {
		panic(err)
	}

	// Init gin and handlers
	r := gin.Default()
	handlers := handlers.InitDB(db) 

	// Get todos
	r.GET("/todos", handlers.GetTodos)
	
	// Get todo by id
	r.GET("/todos/:id", handlers.GetTodoById)

	// Add todo
	r.POST("/todos", handlers.AddTodo)
	
	// Update todo by id
	r.PATCH("/todos/:id", handlers.UpdateTodoById)

	// Delete todo by id
	r.DELETE("/todos/:id", handlers.DeleteTodoById)

	r.Run()
}
