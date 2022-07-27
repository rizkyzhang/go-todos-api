package handlers

import (
	"fmt"
	"go-todos-api/models"
	"go-todos-api/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h* Handler) DeleteTodoById(ctx *gin.Context) {
		var deletedTodo models.Todo

		// Get id param and convert it to int
		id := ctx.Param("id")
		intId, err := strconv.Atoi(id)

		if (err != nil) {
			return
		}

		// Validate id param
		rowCount := 0
		sqlStatement := `SELECT COUNT(*) as count FROM todos;`

		err = h.db.QueryRow(sqlStatement).Scan(&rowCount)

		if (err != nil) {
			panic(err)
		}

		if (intId < 1 || intId > rowCount) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Todo %s not found", id),
				"message": fmt.Sprintf("Delete todo %s failed", id),
				"status": http.StatusNotFound,
				"todos": nil,
				"deleted_todo": nil,
			})

			return 
		}

		// Delete todo
		sqlStatement = `
		DELETE FROM todos WHERE id = $1
		RETURNING id, is_completed, todo;
		`
		err = h.db.QueryRow(sqlStatement, intId).Scan(&deletedTodo.Id, &deletedTodo.IsCompleted, &deletedTodo.Todo)

		if (err != nil) {
			panic(err)
		}

		// Get todos
		todos := utils.GetTodosDB(h.db)

		ctx.JSON(http.StatusOK, gin.H{
			"deleted_todo": deletedTodo,
			"message": fmt.Sprintf("Delete todo %s success", id),
			"status": http.StatusOK,
			"todos": todos,
		})
}