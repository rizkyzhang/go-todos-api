package handlers

import (
	"net/http"

	"go-todos-api/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func (h *Handler) GetTodos(ctx *gin.Context) {
	todos := utils.GetTodosDB(h.db)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Get todos success",
		"status": http.StatusOK,
		"todos": todos,
	})
}