package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

//go:embed tools/*.json
var toolsFS embed.FS

func main() {
	var prompt string
	flag.StringVar(&prompt, "p", "", "Prompt to send to LLM")
	flag.Parse()

	if prompt == "" {
		panic("Prompt must not be empty")
	}

	apiKey := os.Getenv("OPENROUTER_API_KEY")
	baseUrl := os.Getenv("OPENROUTER_BASE_URL")
	if baseUrl == "" {
		baseUrl = "https://openrouter.ai/api/v1"
	}

	if apiKey == "" {
		panic("Env variable OPENROUTER_API_KEY not found")
	}

	toolFiles, err := toolsFS.ReadDir("tools")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading embedded tools: %v\n", err)
		os.Exit(1)
	}

	var tools []openai.ChatCompletionToolUnionParam
	for _, toolFile := range toolFiles {
		data, err := toolsFS.ReadFile("tools/" + toolFile.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading %s: %v\n", toolFile.Name(), err)
			os.Exit(1)
		}

		var tool openai.ChatCompletionToolUnionParam
		if err := json.Unmarshal(data, &tool); err != nil {
			fmt.Fprintf(os.Stderr, "error unmarshalling %s: %v\n", toolFile.Name(), err)
			os.Exit(1)
		}

		tools = append(tools, tool)
	}

	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL(baseUrl))

	messages := []openai.ChatCompletionMessageParamUnion{
		openai.UserMessage(prompt),
	}

	for {
		resp, err := client.Chat.Completions.New(context.Background(),
			openai.ChatCompletionNewParams{
				Model:    "anthropic/claude-haiku-4.5",
				Messages: messages,
				Tools:    tools,
			},
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		if len(resp.Choices) == 0 {
			panic("No choices in response")
		}

		choice := resp.Choices[0]

		if len(choice.Message.ToolCalls) == 0 {
			fmt.Print(choice.Message.Content)
			break
		}

		toolCallParams := make([]openai.ChatCompletionMessageToolCallUnionParam, len(choice.Message.ToolCalls))
		for i, tc := range choice.Message.ToolCalls {
			toolCallParams[i] = tc.ToParam()
		}

		messages = append(messages, openai.ChatCompletionMessageParamUnion{
			OfAssistant: &openai.ChatCompletionAssistantMessageParam{
				ToolCalls: toolCallParams,
			},
		})

		for _, toolCall := range choice.Message.ToolCalls {
			fmt.Fprintf(os.Stderr, "Tool call: %s(%s)\n", toolCall.Function.Name, toolCall.Function.Arguments)
			result := executeToolCall(toolCall.Function.Name, toolCall.Function.Arguments)
			messages = append(messages, openai.ToolMessage(result, toolCall.ID))
		}
	}
}

func executeToolCall(name, arguments string) string {
	switch name {
	case "Read":
		return executeReadTool(arguments)
	case "Write":
		return executeWriteTool(arguments)
	case "Bash":
		return executeBashTool(arguments)
	default:
		return fmt.Sprintf("unknown tool: %s", name)
	}
}

func executeReadTool(arguments string) string {
	var args struct {
		FilePath string `json:"file_path"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("error parsing arguments: %v", err)
	}

	data, err := os.ReadFile(args.FilePath)
	if err != nil {
		return fmt.Sprintf("error reading file: %v", err)
	}

	return string(data)
}

func executeWriteTool(arguments string) string {
	var args struct {
		FilePath string `json:"file_path"`
		Content  string `json:"content"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("error parsing arguments: %v", err)
	}

	file, err := os.Create(args.FilePath)
	if err != nil {
		return fmt.Sprintf("error creating or truncating file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(args.Content); err != nil {
		return fmt.Sprintf("error writing to file: %v", err)
	}

	return fmt.Sprintf("successfully written to: %s", args.FilePath)
}

func executeBashTool(arguments string) string {
	var args struct {
		Command string `json:"command"`
	}

	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf("error parsing arguments: %v", err)
	}

	cmd := exec.Command("sh", "-c", args.Command)

	out, _ := cmd.CombinedOutput()

	return string(out)
}