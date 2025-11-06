package pipeline

import (
	"errors"

	"github.com/pretodev/agent_pipeline/internal/task"
)

var (
	ErrFileNotFound    = errors.New("pipeline file not found")
	ErrInvalidYAML     = errors.New("invalid YAML format")
	ErrEmptyTasks      = errors.New("pipeline must have at least one task")
	ErrMissingTaskName = errors.New("task name is required")
)

type Pipeline struct {
	Description string
	Tasks       []task.Task
}
