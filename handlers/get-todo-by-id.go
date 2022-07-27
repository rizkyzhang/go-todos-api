package handlers

import (
	"fmt"
	"go-todos-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetTodoById(ctx *gin.Context) {
	var todo models.Todo

	// Get id param and convert it to int
	id := ctx.Param("id") 
	intId, err := strconv.Atoi(id)

	rowCount := 0

	if (err != nil) {
		panic(err) 
	}

	// Validate id param
	err = h.db.QueryRow(`Select COUNT(*) as count FROM todos;`).Scan(&rowCount)

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

	err = h.db.QueryRow(`SELECT * FROM todos WHERE id = $1;`, intId).Scan(&todo.Id, &todo.IsCompleted, &todo.Todo)

	if (err != nil) {
		panic(err) 
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Get todo %s success", id),
		"status": http.StatusOK,
		"todo": todo,
	})
}