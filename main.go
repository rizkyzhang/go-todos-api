package main

import (
	"net/http"
	"strconv"

	"fmt"

	"github.com/gin-gonic/gin"
)

type todo struct {
	Id uint `json:"id"`
	Todo string `json:"todo"`
	IsCompleted bool `json:"is_completed"`
}

func main() {
	todos := []todo{
		{Id: 0, Todo: "Test", IsCompleted: false},
		{Id: 1, Todo: "Push commit to github", IsCompleted: false},
		{Id: 2, Todo: "Clean code", IsCompleted: false},
	}

	r := gin.Default()

	// Get todos
	r.GET("/todos", func (ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Get todos success",
			"status": http.StatusOK,
			"todos": todos,
		})
	})
	
	// Get todo by id
	r.GET("/todos/:id", func(ctx *gin.Context) {
		id := ctx.Param("id") 
		intId, err := strconv.Atoi(id)

		if (err != nil) {
			return 
		}

		if (intId < 0 || intId > len(todos) - 1) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Todo %s not found", id),
				"message": fmt.Sprintf("Get todo %s failed", id),
				"status": http.StatusNotFound,
				"todo": nil,
			})

			return 
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Get todo %s success", id),
			"status": http.StatusOK,
			"todo": todos[intId],
		})
	})

	// Add todo
	r.POST("/todos", func(ctx *gin.Context) {
		var newTodo todo

		if err := ctx.BindJSON(&newTodo); err != nil {
			return
		}

		todos = append(todos, newTodo)

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Add todo success",
			"newTodo": newTodo,
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

	r.Run()
}