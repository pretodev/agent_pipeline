package main

import (
	"fmt"
	"os"

	"github.com/pretodev/agent_pipeline/internal/pipeline"
	"github.com/pretodev/agent_pipeline/internal/task/script"
)

func main() {
	p, err := pipeline.ParseFile("./example/helloworld.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
		return
	}

	if len(p.Tasks) == 0 {
		os.Exit(0)
		return
	}

	for _, task := range p.Tasks {
		fmt.Println(task.Name)
		executor := script.Bash()
		if err := executor.Parse(task); err != nil {
			fmt.Println(err)
			os.Exit(1)
			return
		}
		fmt.Println(executor)
	}
}
