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
	r.GET("/todos", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Get todos success",
			"status": http.StatusOK,
			"todos": todos,
		})
	})
	
	// Get todo by id
	r.GET("/todos/:id", func(c *gin.Context) {
		id := c.Param("id") 
		intId, err := strconv.Atoi(id)

		if (err != nil) {
			return 
		}

		if (intId < 0 || intId > len(todos) - 1) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": fmt.Sprintf("Todo %s not found", id),
				"message": fmt.Sprintf("Get todo %s failed", id),
				"status": http.StatusNotFound,
				"todo": nil,
			})

			return 
		}

		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Get todo %s success", id),
			"status": http.StatusOK,
			"todo": todos[intId],
		})
	})

	// Add todo
	r.POST("/todos", func(c *gin.Context) {
		var newTodo todo

		if err := c.BindJSON(&newTodo); err != nil {
			return
		}

		todos = append(todos, newTodo)

		c.JSON(http.StatusOK, gin.H{
			"message": "Add todo success",
			"newTodo": newTodo,
			"status": http.StatusOK,
			"todos": todos,
		})
	})

	r.Run()
}