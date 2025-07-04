package main

import (
	"fmt"
	"log"

	"todo_watcher/model"

	"github.com/gin-gonic/gin"
)

func getTasks(path string) func(c *gin.Context) {
	return func(c *gin.Context) {
		todos, err := model.ParseFiles(path)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, todos)
	}
}

func main() {
	path := "C:/Users/richa/Desktop/todo"
	todos, err := model.ParseFiles(path)
	if err != nil {
		log.Fatalf("Error reading files: %v", err)
	}
	for _, todo := range todos {
		fmt.Println(todo)
	}

	router := gin.Default()
	router.GET("/tasks", getTasks(path))
	router.Run("localhost:8080")
}
