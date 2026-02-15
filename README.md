# Gobot 
Gobot is a high-performance, CLI-based AI coding agent built in Go. Inspired by Claude Code, it is designed to help developers navigate, understand, and edit complex codebases directly from the terminal with near-zero latency.
---
Unlike standard chat-based AI interfaces, Gobot is built to be a terminal-first pair programmer that understands your local file system, runs commands, and performs multi-file edits autonomously.
## Project Structure

```text
.
├── app/
│   └── tools/
│       └── read.json
│       └── write.json                
│       └── bash.json                                
│   └── main.go        
├── go.mod              
├── go.sum
└── README.md
```

## Features
- Terminal-Native: Minimalist and fast CLI interface built with Go for high performance.
- Context-Aware: Automatically indexes and searches your local codebase to provide relevant answers.
- File System Operations: Capabilities to read, create, and edit files across your entire project.
- Tool Use (Function Calling): Integrated support for executing shell commands and interpreting compiler/test feedback.
- Built for Speed: Leverages Go’s concurrency model to handle file processing and API requests efficiently.

## Tech Stack
- Language: Go (Golang)
- AI Models: Supports Anthropic (claude-haiku-4.5) & OpenAI API
- Performance: Optimized for large-scale codebases and rapid context retrieval.

## Installation

Since Gobot is built in Go, you can install it as a single static binary.

### Clone the repository
```
git clone https://github.com/madden-go/Gobot.git
cd Gobot
```

### Build the project
```
go build -o gobot ./app/main.go
```

### (Optional) Move to your path
```
sudo mv gobot /usr/local/bin/
```

## Configuration

Set your API keys as environment variables:

```
export ANTHROPIC_API_KEY="your-api-key"
```
#### or
```
export OPENAI_API_KEY="your-api-key"
```

## Usage

Start the agent in your project root:
```
gobot
```

## Why Go?

While many agents are written in Python, Gobot uses Go to provide a single binary that starts instantly, uses minimal RAM, and handles concurrent I/O tasks efficiently—making it perfect for developers on Arch Linux or those working in performance-critical environments.
