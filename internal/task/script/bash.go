package script

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/pretodev/agent_pipeline/internal/task"
)

const (
	pwdKey      = "pwd"
	commandsKey = "commands"
	sourceKey   = "source"
)

var (
	ErrInvalidPwdType      = errors.New("pwd must be a string")
	ErrPwdNotFound         = errors.New("pwd directory does not exist")
	ErrInvalidCommandsType = errors.New("commands must be a list of strings")
	ErrCommandNotFound     = errors.New("command not found in PATH")
	ErrInvalidSourceType   = errors.New("source must be a string")
	ErrSourceRequired      = errors.New("source is required")
)

type BashTaskExecutor struct {
	pwd      string
	commands []string
	source   string
}

func (t *BashTaskExecutor) Execute() error {
	return nil
}

func (t *BashTaskExecutor) Parse(task task.Task) error {
	if err := t.parsePwd(task.Options); err != nil {
		return err
	}

	if err := t.parseCommands(task.Options); err != nil {
		return err
	}

	if err := t.parseSource(task.Options); err != nil {
		return err
	}

	return nil
}

func (t *BashTaskExecutor) parsePwd(options map[string]any) error {
	if pwd, ok := options[pwdKey]; ok {
		pwdStr, ok := pwd.(string)
		if !ok {
			return ErrInvalidPwdType
		}

		if pwdStr == "" {
			currentDir, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			t.pwd = currentDir
			return nil
		}

		if _, err := os.Stat(pwdStr); os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", ErrPwdNotFound, pwdStr)
		}

		t.pwd = pwdStr
		return nil
	}

	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	t.pwd = currentDir
	return nil
}

func (t *BashTaskExecutor) parseCommands(options map[string]any) error {
	if commands, ok := options[commandsKey]; ok {
		commandsList, ok := commands.([]interface{})
		if !ok {
			return ErrInvalidCommandsType
		}

		t.commands = make([]string, 0, len(commandsList))
		for i, cmd := range commandsList {
			cmdStr, ok := cmd.(string)
			if !ok {
				return fmt.Errorf("%w: command at index %d is not a string", ErrInvalidCommandsType, i)
			}

			if _, err := exec.LookPath(cmdStr); err != nil {
				return fmt.Errorf("%w: %s", ErrCommandNotFound, cmdStr)
			}

			t.commands = append(t.commands, cmdStr)
		}
	}

	return nil
}

func (t *BashTaskExecutor) parseSource(options map[string]any) error {
	source, ok := options[sourceKey]
	if !ok {
		return ErrSourceRequired
	}

	sourceStr, ok := source.(string)
	if !ok {
		return ErrInvalidSourceType
	}

	t.source = sourceStr
	return nil
}
