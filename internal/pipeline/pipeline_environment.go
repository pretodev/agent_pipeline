package pipeline

import (
	"os"
	"regexp"
	"strings"
)

type Environment struct {
	variables map[string]string
}

func NewEnvironment() *Environment {
	return &Environment{
		variables: make(map[string]string),
	}
}

func (e *Environment) LoadFromOS() {
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			e.variables[parts[0]] = parts[1]
		}
	}
}

func (e *Environment) Get(key string) string {
	return e.variables[key]
}

func (e *Environment) Set(key, value string) {
	e.variables[key] = value
}

func (e *Environment) GetAll() map[string]string {
	result := make(map[string]string, len(e.variables))
	for k, v := range e.variables {
		result[k] = v
	}
	return result
}

func (e *Environment) ExpandString(s string) string {
	varPattern := regexp.MustCompile(`\$\{([^}]+)\}`)

	return varPattern.ReplaceAllStringFunc(s, func(match string) string {
		varName := strings.TrimPrefix(strings.TrimSuffix(match, "}"), "${")
		if value, exists := e.variables[varName]; exists {
			return value
		}
		return ""
	})
}

func (e *Environment) ExpandValue(value interface{}) interface{} {
	switch v := value.(type) {
	case string:
		return e.ExpandString(v)
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = e.ExpandValue(item)
		}
		return result
	case map[string]interface{}:
		result := make(map[string]interface{}, len(v))
		for key, val := range v {
			result[key] = e.ExpandValue(val)
		}
		return result
	default:
		return value
	}
}
