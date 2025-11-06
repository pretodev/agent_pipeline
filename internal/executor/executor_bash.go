package executor

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pretodev/agent_pipeline/internal/pipeline"
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
	ErrScriptExecution     = errors.New("script execution failed")
	ErrNotParsed           = errors.New("executor not parsed, call Parse before Execute")
)

type bashexec struct {
	pwd      string
	commands []string
	source   string
}

func Bash() *bashexec {
	return &bashexec{}
}

func (t *bashexec) String() string {
	return fmt.Sprintf("PWD: %s\n", t.pwd)
}

func (t *bashexec) Execute(env *pipeline.Environment) error {
	if t.source == "" {
		return ErrNotParsed
	}

	cmd := exec.Command("bash", "-c", t.source)
	cmd.Dir = t.pwd

	if env != nil {
		envVars := os.Environ()
		for key, value := range env.GetAll() {
			envVars = append(envVars, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = envVars
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf("%w: %s", ErrScriptExecution, stderr.String())
		}
		return fmt.Errorf("%w: %v", ErrScriptExecution, err)
	}

	if stdout.Len() > 0 {
		t.processOutput(stdout.String(), env)
	}

	return nil
}

func (t *bashexec) processOutput(output string, env *pipeline.Environment) {
	setVarPattern := regexp.MustCompile(`^#\[setVariable '([^']+)';](.*)$`)
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if matches := setVarPattern.FindStringSubmatch(line); matches != nil {
			varName := matches[1]
			varValue := matches[2]
			if env != nil {
				env.Set(varName, varValue)
			}
			continue
		}
		fmt.Println(line)
	}
}

func (t *bashexec) Parse(task pipeline.Task) error {
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

func (t *bashexec) parsePwd(options map[string]any) error {
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

func (t *bashexec) parseCommands(options map[string]any) error {
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

func (t *bashexec) parseSource(options map[string]any) error {
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
