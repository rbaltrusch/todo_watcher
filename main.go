package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	NOT_STARTED = iota
	IN_PROGRESS = iota
	COMPLETED   = iota
)

const (
	DIVIDER_GROUP  = iota
	INDENTED_GROUP = iota
	UNGROUPED      = iota
)

type Todo struct {
	Source   string `json:"source"`
	Content  string `json:"content"`
	Status   int    `json:"status"`
	Subtasks []Todo `json:"subtasks,omitempty"`
}

func (t Todo) formatStatus() string {
	switch t.Status {
	case NOT_STARTED:
		return "Not Started"
	case IN_PROGRESS:
		return "In Progress"
	case COMPLETED:
		return "Completed"
	default:
		return "Unknown"
	}
}

func (t Todo) String() string {
	status := t.formatStatus()
	headline := fmt.Sprintf("Todo (group=%t): %s >>> %s. Status: %s", t.HasSubtasks(), t.Source, t.Content, status)
	if t.HasSubtasks() {
		return headline
	}

	subtaskStrings := make([]string, len(t.Subtasks))
	for i, subtask := range t.Subtasks {
		subtaskStrings[i] = "\t" + subtask.String()
	}
	return fmt.Sprintf("%s\nSubtasks:\n%s", headline, strings.Join(subtaskStrings, "\n"))
}

func (t Todo) HasSubtasks() bool {
	return len(t.Subtasks) > 0
}

func determineStatus(status string) int {
	fmt.Println("Determining status for:", status)
	switch strings.TrimSpace(status) {
	case "x":
		return COMPLETED
	case "~":
		return IN_PROGRESS
	default:
		return NOT_STARTED
	}
}

func readFiles(path string) ([]Todo, error) {
	var todos []Todo
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`^\s*([x~])?\s*(.*)$`)
	dividerPattern := regexp.MustCompile(`^\s*([-=]+)\s*(.*)$`)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := path + "/" + entry.Name()
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var group *Todo
		var isDividerGroup bool
		groupType := UNGROUPED
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "" {
				continue
			}

			// skip simple divider lines, but not labelled dividers
			dividerMatches := dividerPattern.FindStringSubmatch(line)
			if length := len(dividerMatches); length > 0 && length <= 2 {
				continue
			}
			isDividerGroup = len(dividerMatches) > 2
			if isDividerGroup {
				groupType = DIVIDER_GROUP
			}
			fmt.Println(line, len(dividerMatches), isDividerGroup, groupType)

			// check for subtasks grouped by indentation
			if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {
				if groupType != INDENTED_GROUP {
					groupType = INDENTED_GROUP
					group = &todos[len(todos)-1]
				}
			} else if groupType == INDENTED_GROUP {
				groupType = UNGROUPED
			}

			matches := re.FindStringSubmatch(line)
			var todo Todo
			if len(matches) < 3 {
				todo.Content = line
				todo.Status = NOT_STARTED
			} else {
				todo.Content = matches[2]
				todo.Status = determineStatus(matches[1])
			}
			if isDividerGroup {
				todo.Content = dividerMatches[2]
				if strings.Contains(todo.Content, "done") {
					todo.Status = COMPLETED
				}
			}
			todo.Source = filePath

			fmt.Println(line, len(dividerMatches), isDividerGroup, groupType, groupType == UNGROUPED, isDividerGroup)
			if groupType == UNGROUPED || isDividerGroup {
				todos = append(todos, todo)
				if isDividerGroup {
					group = &todos[len(todos)-1]
				}
			} else {
				if group.Status == COMPLETED {
					todo.Status = COMPLETED
				}
				group.Subtasks = append(group.Subtasks, todo)
			}
		}

		err = scanner.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	return todos, nil
}

func getTasks(path string) func(c *gin.Context) {
	return func(c *gin.Context) {
		todos, err := readFiles(path)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, todos)
	}
}

func main() {
	path := "C:/Users/richa/Desktop/todo"
	router := gin.Default()
	router.GET("/tasks", getTasks(path))
	router.Run("localhost:8080")
}
