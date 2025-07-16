package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"

	"todo_watcher/filewatcher"
	"todo_watcher/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
func openFile(todoPath string, editor string) func(c *gin.Context) {
	return func(c *gin.Context) {
		filePath := c.Query("file")
		if filePath == "" {
			c.JSON(400, gin.H{"error": "file query parameter is required"})
			return
		}

		execPath, err := exec.LookPath(editor)
		if err != nil {
			c.JSON(500, gin.H{"error": "could not find editor.", "editor": editor})
			return
		}

		// Check if the file exists
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			c.JSON(404, gin.H{"error": fmt.Sprintf("file not found: %s", filePath)})
			return
		}

		// Check if the file exists in the todo folder
		if path.Dir(filePath) != todoPath {
			c.JSON(400, gin.H{"error": fmt.Sprintf("file %s is not in the todo folder %s", filePath, todoPath)})
			return
		}

		cmd := exec.Command(execPath, filePath)
		cmd.Dir = todoPath
		err = cmd.Start()
		if err != nil {
			c.JSON(500, gin.H{"error": fmt.Sprintf("could not open file: %v", err)})
			return
		}

		c.JSON(200, gin.H{"status": "OK"})
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all connections by default
	},
}

func makeWebSocketHandler(broadcast chan string) gin.HandlerFunc {
	return func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println("WebSocket Upgrade failed:", err)
			return
		}
		defer conn.Close()

		for {
			// Write a message to the WebSocket connection
			message := <-broadcast
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				fmt.Println("Write error:", err)
				break
			}
		}
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
	editor := getEnvOrDefault("EDITOR", "code")

	broadcast := make(chan string)
	watcher := filewatcher.CreateWatcher(path)
	defer watcher.Close()
	go filewatcher.HandleFileEvents(watcher, broadcast)

	router := gin.Default()
	router.GET("/ws", makeWebSocketHandler(broadcast))
	router.GET("/api/todos", getTasks(path))
	router.GET("/api/open", openFile(path, editor))
	router.Run(address)
}
