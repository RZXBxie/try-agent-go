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

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Please set GEMINI_API_KEY environment variable.")
		os.Exit(1)
	}

	client := NewGeminiClient(apiKey, "gemini-2.5-flash")

	// 切换提示词风格：修改这里的参数即可
	// 可选: personal-assistant | sarcastic-friend | coding-mentor | anime-girl | strict-engineer
	systemPrompt, err := GetSystemPrompt("personal-assistant")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load system prompt: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Gemini Chat (type '/exit' to quit)\n")

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		if strings.TrimSpace(strings.ToLower(input)) == "/exit" {
			break
		}

		if strings.TrimSpace(input) == "" {
			continue
		}

		reply, err := client.SendMessage(input, systemPrompt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError: %v\n\n", err)
			continue
		}

		fmt.Printf("\nGemini: %s\n\n", reply)
	}
}
