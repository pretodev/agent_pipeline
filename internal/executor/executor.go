package executor

import "github.com/pretodev/agent_pipeline/internal/pipeline"

type Task interface {
	Parse(task pipeline.Task) error
	Execute(env *pipeline.Environment) error
	String() string
}
