package main

import (
	"fmt"
	"log"

	"todo_watcher/model"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	godotenv.Load()
	path := os.Getenv("TODO_FOLDER")
	if path == "" {
		log.Fatal("TODO_FOLDER environment variable must be set")
	}

	address := fmt.Sprintf("%s:%s", getEnvOrDefault("HOST", "localhost"), getEnvOrDefault("PORT", "8080"))

	router := gin.Default()
	router.GET("/api/todos", getTasks(path))
	router.Run(address)
}
