package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type todo struct {
	Id uint `json:"id"`
	Todo string `json:"todo" binding:"required"`
	IsCompleted bool `json:"is_completed"`
}

func getTodosDB(db *sql.DB) ([]todo) {
	var todos []todo

	sqlStatement := `SELECT * FROM todos`
	rows, err := db.Query(sqlStatement)

	if (err != nil) {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		var todo todo

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

func main() {
	const (
		host = "localhost"
		port = 5432
		user = "postgres"
		password = "root"
		dbname = "postgres"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	defer db.Close()

	todos := []todo{
		{Id: 0, Todo: "Test", IsCompleted: false},
		{Id: 1, Todo: "Push commit to github", IsCompleted: false},
		{Id: 2, Todo: "Clean code", IsCompleted: false},
	}

	r := gin.Default()

	// Get todos
	r.GET("/todos", func (ctx *gin.Context) {
		todos := getTodosDB(db)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Get todos success",
			"status": http.StatusOK,
			"todos": todos,
		})
	})
	
	// Get todo by id
	r.GET("/todos/:id", func(ctx *gin.Context) {
		var todo todo

		id := ctx.Param("id") 
		intId, err := strconv.Atoi(id)

		rowCount := 0

		if (err != nil) {
			panic(err) 
		}

		err = db.QueryRow(`Select COUNT(*) as count FROM todos;`).Scan(&rowCount)

		if (err != nil) {
			panic(err)
		}

		if (intId < 1 || intId > rowCount) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Todo %s not found", id),
				"message": fmt.Sprintf("Get todo %s failed", id),
				"status": http.StatusNotFound,
				"todo": nil,
			})

			return 
		}

		err = db.QueryRow(`SELECT * FROM todos WHERE id = $1;`, intId).Scan(&todo.Id, &todo.IsCompleted, &todo.Todo)

		if (err != nil) {
			panic(err) 
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Get todo %s success", id),
			"status": http.StatusOK,
			"todo": todo,
		})
	})

	// Add todo
	r.POST("/todos", func(ctx *gin.Context) {
		var reqBody todo
		var newTodo todo

		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "JSON Payload is not valid",
				"message": "Add todo failed",
				"new_todo": nil,
				"status": http.StatusBadRequest,
				"todos": nil,
			})

			return
		}

		sqlStatement := `
		INSERT INTO todos (todo)
		VALUES ($1)
		RETURNING id, is_completed, todo;
		`
		err = db.QueryRow(sqlStatement, reqBody.Todo).Scan(&newTodo.Id, &newTodo.IsCompleted, &newTodo.Todo)

		if (err != nil) {
			panic(err)
		}

		todos := getTodosDB(db)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Add todo success",
			"new_todo": newTodo,
			"status": http.StatusOK,
			"todos": todos,
		})
	})

	// Update todo by id
	r.PATCH("/todos/:id", func(ctx *gin.Context) {
		var updatedTodo todo

		id := ctx.Param("id")
		intId, err := strconv.Atoi(id)

		if (err != nil) {
			return 
		}

		if (intId < 0 || intId > len(todos) - 1) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Todo %s not found", id),
				"message": fmt.Sprintf("Update todo %s failed", id),
				"status": http.StatusNotFound,
				"todos": nil,
				"updated_todo": nil,
			})

			return 
		}

		if err := ctx.BindJSON(&updatedTodo); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "JSON Payload is not valid",
				"message": fmt.Sprintf("Update todo %s failed", id),
				"status": http.StatusBadRequest,
				"todos": nil,
				"updated_todo": nil,
			})

			return
		}

		todos[intId] = updatedTodo
		todos[intId].Id = uint(intId)

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Update todo %s success", id),
			"status": http.StatusOK,
			"todos": todos,
			"updated_todo": todos[intId],
		})
	})

	// Delete todo by id
	r.DELETE("/todos/:id", func(ctx *gin.Context) {
		var updatedTodos []todo

		id := ctx.Param("id")
		intId, err := strconv.Atoi(id)

		if (err != nil) {
			return
		}

		if (intId < 0 || intId > len(todos) - 1) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Todo %s not found", id),
				"message": fmt.Sprintf("Delete todo %s failed", id),
				"status": http.StatusNotFound,
				"todos": nil,
				"deleted_todo": nil,
			})

			return 
		}

		deletedTodo := todos[intId]

		for i, todo := range todos {
			if i != intId {
				updatedTodos = append(updatedTodos, todo)
			}
		}

		todos = updatedTodos

		ctx.JSON(http.StatusOK, gin.H{
			"deleted_todo": deletedTodo,
			"message": fmt.Sprintf("Delete todo %s success", id),
			"status": http.StatusOK,
			"todos": updatedTodos,
		})
	})

	r.Run()
}