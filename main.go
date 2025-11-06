package main

import (
	"fmt"
	"os"

	"github.com/pretodev/agent_pipeline/internal/executor"
	"github.com/pretodev/agent_pipeline/internal/pipeline"
)

func main() {
	p, err := pipeline.ParseFile("./example/bash.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing pipeline: %v\n", err)
		os.Exit(1)
	}

	if len(p.Tasks) == 0 {
		return
	}

	for _, task := range p.Tasks {
		fmt.Println(task.Name)

		executor := executor.Bash()

		if err := executor.Parse(task); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing task: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(executor)

		if err := executor.Execute(p.Environment); err != nil {
			fmt.Fprintf(os.Stderr, "Error executing task: %v\n", err)
			os.Exit(1)
		}
	}
}
