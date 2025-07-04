package model

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
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
	Source   string
	Content  string
	Status   int
	Subtasks []Todo
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

type Result struct {
	Scanner  *bufio.Scanner
	Err      error
	Filepath string
	Close    func()
}

func createFileReaderIterator(path string) (<-chan Result, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	ch := make(chan Result)
	go func() {
		defer close(ch)

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filePath := path + "/" + entry.Name()
			file, err := os.Open(filePath)
			if err != nil {
				ch <- Result{Err: err}
				continue
			}

			scanner := bufio.NewScanner(file)
			ch <- Result{Scanner: scanner, Filepath: filePath, Close: func() { file.Close() }}
		}
	}()
	return ch, nil
}

type TodoParseResult struct {
	Todos     []Todo
	Group     *Todo
	GroupType int
}

var stripPattern = regexp.MustCompile(`\s*:\s*$`)
var pattern = regexp.MustCompile(`^\s*([x~])?\s*(.*)$`)
var dividerPattern = regexp.MustCompile(`^\s*([-=]+)\s*(.*)$`)

func parseTodoLine(line string, filepath string, result *TodoParseResult) {
	group := result.Group
	groupType := result.GroupType

	if strings.TrimSpace(line) == "" {
		return
	}

	// skip simple divider lines, but not labelled dividers
	dividerMatches := dividerPattern.FindStringSubmatch(line)
	if length := len(dividerMatches); length > 0 && length <= 2 {
		return
	}
	isDivider := len(dividerMatches) > 2
	if isDivider {
		groupType = DIVIDER_GROUP
	}
	// fmt.Println(line, len(dividerMatches), isDividerGroup, groupType)

	// check for subtasks grouped by indentation
	if strings.HasPrefix(line, "\t") || strings.HasPrefix(line, " ") {
		if groupType != INDENTED_GROUP {
			groupType = INDENTED_GROUP
			group = &result.Todos[len(result.Todos)-1]
		}
	} else if groupType == INDENTED_GROUP {
		groupType = UNGROUPED
	}

	matches := pattern.FindStringSubmatch(line)
	var todo Todo
	if len(matches) < 3 {
		todo.Content = line
		todo.Status = NOT_STARTED
	} else {
		todo.Content = matches[2]
		todo.Status = determineStatus(matches[1])
	}
	if isDivider {
		todo.Content = dividerMatches[2]
		if strings.Contains(todo.Content, "done") {
			todo.Status = COMPLETED
		}
	}
	todo.Content = stripPattern.ReplaceAllString(todo.Content, "")
	todo.Source = filepath

	// fmt.Println(line, len(dividerMatches), isDividerGroup, groupType, groupType == UNGROUPED)
	if groupType == UNGROUPED || groupType == DIVIDER_GROUP {
		result.Todos = append(result.Todos, todo)
		if groupType == DIVIDER_GROUP {
			group = &result.Todos[len(result.Todos)-1]
		}
	} else {
		if group.Status == COMPLETED {
			todo.Status = COMPLETED
		}
		group.Subtasks = append(group.Subtasks, todo)
	}

	result.Group = group
	result.GroupType = groupType
}

func ParseFiles(path string) ([]Todo, error) {
	parseResult := TodoParseResult{}
	iter, err := createFileReaderIterator(path)
	if err != nil {
		return nil, err
	}

	for file := range iter {
		if file.Err != nil {
			log.Printf("Error reading file: %v", file.Err)
			continue
		}

		scanner := file.Scanner
		for scanner.Scan() {
			line := scanner.Text()
			parseTodoLine(line, file.Filepath, &parseResult)
		}
		file.Close()
	}
	return parseResult.Todos, nil
}
