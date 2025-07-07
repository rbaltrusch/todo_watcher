package model

import (
	"fmt"
	"strings"
	"time"
)

type Todo struct {
	Status    int        `json:"status"`
	Date      *time.Time `json:"date,omitempty"`
	Source    string     `json:"source,omitempty"`
	Content   string     `json:"content,omitempty"`
	Priority  int        `json:"priority,omitempty"`
	Tentative bool       `json:"tentative,omitempty"`
	SubTasks  []*Todo    `json:"subtasks,omitempty"`
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

func (t Todo) formatPriority() string {
	switch t.Priority {
	case LOW_PRIORITY:
		return "Low"
	case MEDIUM_PRIORITY:
		return "Medium"
	case HIGH_PRIORITY:
		return "High"
	default:
		return "Unknown"
	}
}

func (t Todo) HeadLine() string {
	status := t.formatStatus()
	priority := t.formatPriority()
	return fmt.Sprintf("Todo (group=%t): %s >>> %s. Status: %s Priority: %s", t.HasSubtasks(), t.Source, t.Content, status, priority)
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
