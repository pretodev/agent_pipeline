package task

type Task struct {
	Name        string
	Description string
	Var         string
	Options     map[string]any
}

type TaskExecutor interface {
	Parse(task Task) error
	Execute() error
	String() string
}
