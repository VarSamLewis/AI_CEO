package llm

  import (
  	"context"
  	"os"
  	"github.com/anthropics/anthropic-sdk-go"
  )

  func CallAnthropic(message string) (string, error) {
  	client := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))

  	systemPrompt := os.Getenv("LLM_SYSTEM_PROMPT")
  	if systemPrompt == "" {
  		systemPrompt = "You are a helpful assistant."
  	}

  	ctx := context.Background()
  	resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
  		Model: anthropic.F("claude-3-5-sonnet-20241022"),
  		System: anthropic.F(systemPrompt),
  		Messages: []anthropic.MessageParam{
  			anthropic.NewUserMessage(anthropic.NewTextBlock(message)),
  		},
  		MaxTokens: anthropic.Int(1024),
  	})

  	if err != nil {
  		return "", err
  	}

  	return resp.Content[0].Text, nil
  }
