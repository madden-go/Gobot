package adapter

import (
	"github.com/codecrafters-io/claude-code-starter-go/internal/tools"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
	"github.com/openai/openai-go/v3/shared"
)

// ToOpenAITools 将工具定义转换为 OpenAI 格式
func ToOpenAITools(toolDefs []tools.ToolDefinition) []openai.ChatCompletionToolUnionParam {
	result := make([]openai.ChatCompletionToolUnionParam, 0, len(toolDefs))
	for _, def := range toolDefs {
		result = append(result, openai.ChatCompletionFunctionTool(shared.FunctionDefinitionParam{
			Name:        def.Name,
			Description: param.NewOpt(def.Description),
			Parameters:  shared.FunctionParameters(def.Parameters),
		}))
	}
	return result
}