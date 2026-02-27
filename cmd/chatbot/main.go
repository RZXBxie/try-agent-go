package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"

	"try_agent_go/internal/app"
	"try_agent_go/internal/infra/prompt"
	"try_agent_go/internal/infra/provider"
)

func main() {
	_ = godotenv.Load()

	// 1. 初始化 Provider
	providerName := os.Getenv("MODEL_PROVIDER")
	if providerName == "" {
		providerName = "gemini-2.5-flash"
	}

	p, err := provider.New(providerName, os.Getenv)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create provider %q: %v\n", providerName, err)
		os.Exit(1)
	}

	// 2. 加载 System Prompt
	promptsDir := filepath.Join(projectRoot(), "prompts")
	systemPrompt, err := prompt.Load(promptsDir, "personal-assistant")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load system prompt: %v\n", err)
		os.Exit(1)
	}

	// 3. 创建 ChatService
	chat := app.NewChatService(p, systemPrompt)

	// 4. REPL 交互循环
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Chat [%s] (type '/exit' to quit, '/model <provider>' to switch)\n\n", chat.ProviderName())

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		trimmed := strings.TrimSpace(input)

		if strings.ToLower(trimmed) == "/exit" {
			break
		}

		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(strings.ToLower(trimmed), "/model ") {
			name := strings.TrimSpace(trimmed[7:])
			newProvider, err := provider.New(name, os.Getenv)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nError: %v\n\n", err)
				continue
			}
			chat.SwitchProvider(newProvider)
			fmt.Printf("\nSwitched to provider: %s\n\n", chat.ProviderName())
			continue
		}

		reply, err := chat.Send(trimmed)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError: %v\n\n", err)
			continue
		}

		fmt.Printf("\n%s: %s\n\n", chat.ProviderName(), reply)
	}
}

func projectRoot() string {
	exe, err := os.Executable()
	if err == nil {
		// When running via `go run`, the executable is in a temp dir,
		// so fall back to working directory.
		dir := filepath.Dir(exe)
		if _, err := os.Stat(filepath.Join(dir, "prompts")); err == nil {
			return dir
		}
	}
	// Fall back to current working directory.
	wd, _ := os.Getwd()
	return wd
}
