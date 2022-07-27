package handlers

import (
	"fmt"
	"go-todos-api/models"
	"go-todos-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpdateTodoById (ctx *gin.Context) {
		var reqBody models.Todo
		var currentTodo models.Todo
		var updatedTodo models.Todo

		// Get id param and convert it to int
		id := ctx.Param("id")
		intId, err := strconv.Atoi(id)

		if (err != nil) {
			return 
		}

		// Validate id param
		rowCount := 0
		err = h.db.QueryRow(`Select COUNT(*) as count FROM todos;`).Scan(&rowCount)

		if (err != nil) {
			panic(err)
		}

		if (intId < 0 || intId > rowCount) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("models.Todo %s not found", id),
				"message": fmt.Sprintf("Update models.Todo %s failed", id),
				"status": http.StatusNotFound,
				"todos": nil,
				"updated_todo": nil,
			})

			return 
		}

		// Validate request body
		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "JSON Payload is not valid",
				"message": fmt.Sprintf("Update models.Todo %s failed", id),
				"status": http.StatusBadRequest,
				"todos": nil,
				"updated_todo": nil,
			})
		
			return
		}

		if (reqBody.IsCompleted == nil && reqBody.Todo == nil) {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Either is_completed or models.Todo must be specified",
				"message": fmt.Sprintf("Update models.Todo %s failed", id),
				"status": http.StatusBadRequest,
				"todos": nil,
				"updated_todo": nil,
			})
		
			return
		}

		// Get current todo
		sqlStatement := `SELECT * FROM todos WHERE id = $1`
		err = h.db.QueryRow(sqlStatement, intId).Scan(&currentTodo.Id, &currentTodo.IsCompleted, &currentTodo.Todo)

		if (err != nil) {
			panic(err)
		}

		if (reqBody.IsCompleted == nil) {
			reqBody.IsCompleted = currentTodo.IsCompleted
		} 

		if (reqBody.Todo == nil) {
			reqBody.Todo = currentTodo.Todo
		}

		// Update todo
		sqlStatement = `
		UPDATE todos
		SET is_completed = $2, models.Todo = $3
		WHERE id = $1
		RETURNING id, is_completed, models.Todo;
		`
		err = h.db.QueryRow(sqlStatement, intId, reqBody.IsCompleted, reqBody.Todo).Scan(&updatedTodo.Id, &updatedTodo.IsCompleted, &updatedTodo.Todo)

		if (err != nil) {
			panic(err)
		}

		// Get todos
		todos := utils.GetTodosDB(h.db)

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Update models.Todo %s success", id),
			"status": http.StatusOK,
			"todos": todos,
			"updated_todo": updatedTodo,
		})
}