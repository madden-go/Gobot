package tools

import (
	"encoding/json"
	"fmt"
)

type Manager struct {
	tools map[string]Tool
}

type Tool interface {
	Name() string
	Description() string
	ParametersSchema() map[string]any
	Execute(input json.RawMessage) (string, error)
}

func NewManager() *Manager {
	m := &Manager{
		tools: make(map[string]Tool),
	}

	return m
}

func (m *Manager) Register(tool Tool) {
	m.tools[tool.Name()] = tool
}

func (m *Manager) Get(name string) (Tool, bool) {
	tool, ok := m.tools[name]
	return tool, ok
}

func (m *Manager) List() []Tool {
	result := make([]Tool, 0, len(m.tools))
	for _, tool := range m.tools {
		result = append(result, tool)
	}

	return result //也许可以直接维护全局的一个表
}


func (m *Manager) ExecuteCall(toolName string, input string) (string, error) {
	tool, ok := m.Get(toolName)
	if !ok {
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
	return tool.Execute(json.RawMessage(input))
}