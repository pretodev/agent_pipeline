package pipeline

import (
	"errors"
)

var (
	ErrFileNotFound    = errors.New("pipeline file not found")
	ErrInvalidYAML     = errors.New("invalid YAML format")
	ErrEmptyTasks      = errors.New("pipeline must have at least one task")
	ErrMissingTaskName = errors.New("task name is required")
)

type Pipeline struct {
	Description string
	Tasks       []Task
	Environment *Environment
}

type Task struct {
	Name        string
	Description string
	Var         string
	Options     map[string]any
}
