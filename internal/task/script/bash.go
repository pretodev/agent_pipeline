package script

import "github.com/pretodev/agent_pipeline/internal/task"

const (
	pwdKey      = "pwd"
	commandsKey = "commands"
	sourceKey   = "source"
)

type BashTaskExecutor struct {
	pwd      string
	commands []string
	source   string
}

func (t *BashTaskExecutor) New(task task.Task) (*BashTaskExecutor, error) {
	executor := &BashTaskExecutor{}
	if pwd, ok := task.Options[pwdKey]; ok {
		executor.pwd = pwd.(string)
	}
	if commands, ok := task.Options[commandsKey]; ok {
		executor.commands = commands.([]string)
	}
	if source, ok := task.Options[sourceKey]; ok {
		executor.source = source.(string)
	}
	return executor, nil
}

func (t *BashTaskExecutor) Execute(task task.Task) error {
	return nil
}
