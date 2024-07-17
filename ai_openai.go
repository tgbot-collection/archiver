package main

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
)

func askOpenAI(userID int64) string {
	var config = openai.DefaultConfig("token")
	config.BaseURL = "https://burn.hair"
	var client = openai.NewClientWithConfig(config)

	chats := getChats(userID)
	if len(chats) < 0 {
		return "nothing yet"
	}
	for _, chat := range chats {
		fmt.Println(chat)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "chatgpt error"
	}

	var content = resp.Choices[0].Message.Content
	fmt.Println(content)
	return content
}
