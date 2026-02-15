package tools


type ToolDefinition struct {
	Type     string
	Name     string
	Description string
	Parameters map[string]any
}

func (m *Manager) ToDefinitions() []ToolDefinition {
	result := make([]ToolDefinition, 0, len(m.tools))
	for _, tool := range m.List() {
		result = append(result, ToolDefinition{
			Type:        "function",
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters:  tool.ParametersSchema(),
		})
	}
	return result
}