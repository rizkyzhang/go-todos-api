package handlers

import (
	"go-todos-api/models"
	"go-todos-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h* Handler) AddTodo (ctx *gin.Context) {
		var reqBody models.Todo
		var newTodo models.Todo

		if err := ctx.BindJSON(&reqBody); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "JSON Payload is not valid",
				"message": "Add models.Todo failed",
				"new_todo": nil,
				"status": http.StatusBadRequest,
				"todos": nil,
			})

			return
		}

		sqlStatement := `
		INSERT INTO todos (models.Todo)
		VALUES ($1)
		RETURNING id, is_completed, models.Todo;
		`
		err := h.db.QueryRow(sqlStatement, reqBody.Todo).Scan(&newTodo.Id, &newTodo.IsCompleted, &newTodo.Todo)

		if (err != nil) {
			panic(err)
		}

		todos := utils.GetTodosDB(h.db)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Add models.Todo success",
			"new_todo": newTodo,
			"status": http.StatusOK,
			"todos": todos,
		})
}