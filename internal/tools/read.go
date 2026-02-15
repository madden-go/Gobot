package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ReadInput struct {
	FilePath string `json:"file_path"`
	Offset   string `json:"offset,omitempty"`
	Limit    string `json:"limit,omitempty"`
}

type ReadTool struct{}

func (r *ReadTool) Name() string {
	return "read_file"
}

func (r *ReadTool) Description() string {
	return `Reads a file from local filesystem.

  IMPORTANT:
  - Supports both absolute and relative paths
  - Relative paths are resolved from the current working directory
  - This tool can read images (PNG, JPG, etc.), PDFs, and Jupyter notebooks
  - For images, contents will be presented visually since this is a multimodal LLM
  - This tool can only read files, not directories
  - Returns content with line numbers starting from 1

  Parameters:
  - file_path (required): The path to the file (absolute or relative)
  - offset (optional): The line number to start reading from (default: 0)
  - limit (optional): The number of lines to read (default: read entire file)`
}

func (r *ReadTool) ParametersSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": map[string]any{
				"type":        "string",
				"description": "The path to the file (absolute or relative to current working directory)",
			},
			"offset": map[string]any{
				"type":        "integer",
				"description": "The line number to start reading from (default: 0)",
			},
			"limit": map[string]any{
				"type":        "integer",
				"description": "The number of lines to read",
			},
		},
		"required": []string{"file_path"},
	}
}

func (r *ReadTool) Execute(input json.RawMessage) (string, error) {
	var params ReadInput
	if err := json.Unmarshal(input, &params); err != nil {
		return "", fmt.Errorf("invalid input: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	var fullPath string
	if filepath.IsAbs(params.FilePath) {
		fullPath = params.FilePath
	} else {
		fullPath = filepath.Join(wd, params.FilePath)
	}

	cleanPath := filepath.Clean(fullPath)
	if strings.Contains(cleanPath, "..") {
		return "", fmt.Errorf("file_path cannot contain '..' for security reasons")
	}

	content, err := os.ReadFile(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file does not exist: %s", cleanPath)
		}
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	start, _ := strconv.Atoi(params.Offset)
	if start < 0 {
		start = 0
	}
	if start > len(lines) {
		start = len(lines)
	}

	end := len(lines)
	limit, _ := strconv.Atoi(params.Limit)
	if limit > 0 && start+limit < len(lines) {
		end = start + limit
	}

	var result strings.Builder
	for i := start; i < end; i++ {
		fmt.Fprintf(&result, "%5d\t%s\n", i+1, lines[i])
	}

	return result.String(), nil
}