package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// -----------------------------------------------------------------------
	// 1. 初始化 Provider（通过 MODEL_PROVIDER 环境变量选择，默认 gemini-2.5-flash）
	// -----------------------------------------------------------------------

	providerName := os.Getenv("MODEL_PROVIDER")
	if providerName == "" {
		providerName = "gemini-2.5-flash"
	}

	provider, err := CreateProvider(providerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create provider %q: %v\n", providerName, err)
		os.Exit(1)
	}

	// -----------------------------------------------------------------------
	// 2. 加载 System Prompt
	// -----------------------------------------------------------------------
	// 切换提示词风格：修改这里的参数即可
	// 可选: personal-assistant | sarcastic-friend | coding-mentor | anime-girl | strict-engineer

	systemPrompt, err := GetSystemPrompt("personal-assistant")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load system prompt: %v\n", err)
		os.Exit(1)
	}

	// -----------------------------------------------------------------------
	// 3. 创建 Chat（维护对话历史）
	// -----------------------------------------------------------------------

	chat := NewChat(provider, systemPrompt)

	// -----------------------------------------------------------------------
	// 4. REPL 交互循环
	// -----------------------------------------------------------------------

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Chat [%s] (type '/exit' to quit, '/model <provider>' to switch)\n\n", provider.Name())

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		trimmed := strings.TrimSpace(input)

		// --- 退出 ---
		if strings.ToLower(trimmed) == "/exit" {
			break
		}

		// --- 空行跳过 ---
		if trimmed == "" {
			continue
		}

		// --- /model <provider> 运行时切换（历史重置） ---
		if strings.HasPrefix(strings.ToLower(trimmed), "/model ") {
			name := strings.TrimSpace(trimmed[7:])
			newProvider, err := CreateProvider(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "\nError: %v\n\n", err)
				continue
			}
			provider = newProvider
			chat = NewChat(provider, systemPrompt)
			fmt.Printf("\nSwitched to provider: %s\n\n", provider.Name())
			continue
		}

		// --- 正常对话（通过 Chat 发送，自动维护历史） ---
		reply, err := chat.Send(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError: %v\n\n", err)
			continue
		}

		fmt.Printf("\n%s: %s\n\n", provider.Name(), reply)
	}
}
