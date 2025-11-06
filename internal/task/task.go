package task

type Task struct {
	Name        string
	Description string
	Var         string
	Options     map[string]any
}

type TaskExecutor interface {
	Execute(task Task) error
}
