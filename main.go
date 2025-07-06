package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

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

// currently unsafe, but should be fine because this is local only
func openFile(path string) func(c *gin.Context) {
	return func(c *gin.Context) {
		filePath := c.Query("file")
		if filePath == "" {
			c.JSON(400, gin.H{"error": "file query parameter is required"})
			return
		}

		execPath, err := exec.LookPath("code")
		if err != nil {
			c.JSON(500, gin.H{"error": "could not find code editor"})
			return
		}
		cmd := exec.Command(execPath, filePath)
		cmd.Dir = path
		err = cmd.Start()
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("could not open file: %v", err)})
			return
		}

		c.JSON(200, gin.H{"status": "OK"})
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
	router.GET("/api/open", openFile(path))
	router.Run(address)
}
