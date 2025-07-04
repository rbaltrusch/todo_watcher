package model

import (
	"fmt"
	"strings"
)

type Todo struct {
	Status   int     `json:"status"`
	Source   string  `json:"source,omitempty"`
	Content  string  `json:"content,omitempty"`
	SubTasks []*Todo `json:"subtasks,omitempty"`
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

func (t Todo) HeadLine() string {
	status := t.formatStatus()
	return fmt.Sprintf("Todo (group=%t): %s >>> %s. Status: %s", t.HasSubtasks(), t.Source, t.Content, status)
}
func (t Todo) String() string {
	headline := t.HeadLine()
	if !t.HasSubtasks() {
		return headline
	}

	subtaskStrings := make([]string, len(t.SubTasks))
	for i, subtask := range t.SubTasks {
		subtaskStrings[i] = "\t" + subtask.String()
	}
	return fmt.Sprintf("%s\nSubtasks:\n%s", headline, strings.Join(subtaskStrings, "\n"))
}

func (t Todo) HasSubtasks() bool {
	return len(t.SubTasks) > 0
}
