package main

import (
	"net/http"

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

	r.GET("/todos", func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": http.StatusOK,
			"message": "Get todos success",
			"todos": todos,
		})
	})

	r.Run()
}