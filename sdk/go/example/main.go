package main

import (
	"context"
	"fmt"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
	"github.com/uni-token/core/sdk/go"
)

func loadApiKey() string {
	// TODO: Implement loading API key from storage
	return ""
}

func saveApiKey(apiKey string) {
	// TODO: Implement saving API key to storage
	fmt.Println("API key:", apiKey)
}

func main() {
	result, err := uniToken.RequestUniTokenOpenAI(uniToken.UniTokenOptions{
		AppName:     "MyApp",
		Description: "This is a test application",
		SavedAPIKey: loadApiKey(),
	})
	if err != nil {
		panic(err)
	}
	saveApiKey(result.APIKey)

	if result.APIKey == "" {
		fmt.Println("User did not grant permission for OpenAI token.")
		return
	}

	chatDemo(result.BaseURL, result.APIKey)
}

func chatDemo(baseURL, apiKey string) {
	client := openai.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey(apiKey),
	)
	stream := client.Chat.Completions.NewStreaming(context.TODO(), openai.ChatCompletionNewParams{
		Model: "gpt-4o-mini",
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a concise assistant."),
			openai.UserMessage("Please write a one-sentence bedtime story."),
		},
	})

	if stream.Err() != nil {
		panic(stream.Err())
	}
	for stream.Next() {
		if stream.Err() != nil {
			panic(stream.Err())
		}
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			println(chunk.Choices[0].Delta.Content)
		}
	}
}
