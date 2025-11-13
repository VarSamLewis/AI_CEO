package llm

  import (
  	"context"
  	"os"
  	"github.com/anthropics/anthropic-sdk-go"
  	"github.com/anthropics/anthropic-sdk-go/option"
  )

  func CallAnthropic(message string) (string, error) {
  	client := anthropic.NewClient(option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")))

  	systemPrompt := os.Getenv("LLM_SYSTEM_PROMPT")
  	if systemPrompt == "" {
  		systemPrompt = ""
  	}

  	ctx := context.Background()
  	resp, err := client.Messages.New(ctx, anthropic.MessageNewParams{
  		Model: anthropic.ModelClaudeSonnet4_5_20250929,
		  System: []anthropic.TextBlockParam{
        {Text: systemPrompt	},
    	},  		
			Messages: []anthropic.MessageParam{
  			anthropic.NewUserMessage(anthropic.NewTextBlock(message)),
  		},
  		MaxTokens: 1024,
  	})

  	if err != nil {
  		return "", err
  	}

  	return resp.Content[0].Text, nil
  }
