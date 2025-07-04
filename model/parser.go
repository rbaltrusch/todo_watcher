package model

import (
	"fmt"
	"log"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"todo_watcher/util"
)

const (
	NOT_STARTED = iota
	IN_PROGRESS = iota
	COMPLETED   = iota
)

func determineStatus(status string) int {
	switch strings.TrimSpace(status) {
	case "x":
		return COMPLETED
	case "~":
		return IN_PROGRESS
	default:
		return NOT_STARTED
	}
}

type TodoParseResult struct {
	Todos             []*Todo
	Latest            *Todo
	Groups            util.Stack[*Todo]
	IndentationLevels util.Stack[int]
}

var stripPattern = regexp.MustCompile(`\s*:\s*$`)
var pattern = regexp.MustCompile(`^\s*([x~])?\s*(.*)$`)
var dividerPattern = regexp.MustCompile(`^\s*([-=]+)\s*(.*)$`)

func determineIndentationLevel(line string) int {
	indentation := 0
	for _, char := range line {
		if char == '\t' {
			indentation += 4 // assuming a tab is equivalent to 4 spaces
		} else if char == ' ' {
			indentation++
		} else {
			break
		}
	}
	return indentation
}

// TODO: extract tentative (ending with question mark)

func parseTodoLine(line string, result *TodoParseResult) {
	todo := Todo{Content: "", Status: NOT_STARTED}

	if strings.TrimSpace(line) == "" {
		return
	}

	// skip simple divider lines, but not labelled dividers
	dividerMatches := dividerPattern.FindStringSubmatch(line)
	if length := len(dividerMatches); length > 0 && length <= 2 {
		return
	}
	isDividerGroupHeader := len(dividerMatches) > 2

	// check for subtasks grouped by indentation
	indentation := determineIndentationLevel(line)
	currentIndentation, err := result.IndentationLevels.Top()
	if err == nil && indentation > currentIndentation {
		result.IndentationLevels.Push(indentation)
		group := result.Latest
		result.Groups.Push(group)
	} else {
		for {
			currentIndentation, err = result.IndentationLevels.Top()
			if err != nil || indentation >= currentIndentation {
				break
			}
			result.Groups.Pop()
			result.IndentationLevels.Pop()
		}
	}

	matches := pattern.FindStringSubmatch(line)
	if isDividerGroupHeader {
		defer result.Groups.Push(&todo)
		todo.Content = dividerMatches[2]
		if strings.Contains(todo.Content, "done") {
			todo.Status = COMPLETED
		}
	} else if len(matches) < 3 {
		todo.Content = line
	} else {
		todo.Content = matches[2]
		todo.Status = determineStatus(matches[1])
	}
	todo.Content = stripPattern.ReplaceAllString(todo.Content, "")

	group, _ := result.Groups.Top()
	if group.Status == COMPLETED {
		todo.Status = COMPLETED
	}
	group.SubTasks = append(group.SubTasks, &todo)
	result.Latest = &todo
}

var datePattern = regexp.MustCompile(`^(\d{6})([^\d])`)

func parseFilePathDate(filePath string) *time.Time {
	filename := filepath.Base(filePath)
	matches := datePattern.FindStringSubmatch(filename)
	if len(matches) < 2 {
		return nil
	}

	// Parse YYMMDD date format
	date, err := time.Parse("060102", matches[1])
	if err != nil {
		return nil
	}
	fmt.Println("Parsed date from filename:", date)
	return &date
}

func ParseFiles(path string) ([]*Todo, error) {
	parseResult := TodoParseResult{}
	iter, err := util.CreateFileReaderIterator(path)
	if err != nil {
		return nil, err
	}

	for file := range iter {
		if file.Err != nil {
			log.Printf("Error reading file: %v", file.Err)
			continue
		}

		todo := Todo{Source: file.FilePath, Date: parseFilePathDate(file.FilePath)}
		parseResult.Todos = append(parseResult.Todos, &todo)
		parseResult.IndentationLevels.Push(0)
		parseResult.Groups.Push(&todo)
		parseResult.Latest = &todo

		scanner := file.Scanner
		for scanner.Scan() {
			line := scanner.Text()
			parseTodoLine(line, &parseResult)
		}
		file.Close()
	}
	return parseResult.Todos, nil
}
