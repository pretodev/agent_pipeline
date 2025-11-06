package pipeline

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type pipelineYAML struct {
	Description string                   `yaml:"description"`
	Tasks       []map[string]interface{} `yaml:"tasks"`
}

func ParseFile(filePath string) (*Pipeline, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", ErrFileNotFound, filePath)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var pipelineData pipelineYAML
	if err := yaml.Unmarshal(data, &pipelineData); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidYAML, err.Error())
	}

	if len(pipelineData.Tasks) == 0 {
		return nil, ErrEmptyTasks
	}

	env := NewEnvironment()
	env.LoadFromOS()

	tasks := make([]Task, 0, len(pipelineData.Tasks))
	for i, taskData := range pipelineData.Tasks {
		taskName, ok := taskData["task"].(string)
		if !ok || taskName == "" {
			return nil, fmt.Errorf("%w at index %d", ErrMissingTaskName, i)
		}

		options := make(map[string]any)
		for key, value := range taskData {
			if key != "task" && key != "description" && key != "var" {
				options[key] = env.ExpandValue(value)
			}
		}

		t := Task{
			Name:    taskName,
			Options: options,
		}

		if desc, ok := taskData["description"].(string); ok {
			t.Description = env.ExpandString(desc)
		}

		if varName, ok := taskData["var"].(string); ok {
			t.Var = varName
		}

		tasks = append(tasks, t)
	}

	return &Pipeline{
		Description: pipelineData.Description,
		Tasks:       tasks,
		Environment: env,
	}, nil
}
