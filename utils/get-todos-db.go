package utils

import (
	"database/sql"

	"go-todos-api/models"

	_ "github.com/lib/pq"
)


func GetTodosDB(db *sql.DB) ([]models.Todo) {
	var todos []models.Todo

	sqlStatement := `SELECT * FROM todos`
	rows, err := db.Query(sqlStatement)

	if (err != nil) {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var todo models.Todo

		err = rows.Scan(&todo.Id, &todo.IsCompleted, &todo.Todo)

		if (err != nil) {
			panic(err)
		}

		todos = append(todos, todo)
	}

	err = rows.Err()

	if err != nil {
		panic(err)
	}

	return todos
}