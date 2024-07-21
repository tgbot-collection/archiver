package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
)

func initAI() *openai.Client {
	var config = openai.DefaultConfig(OpenAIKey)
	config.BaseURL = "https://burn.hair/v1"
	var client = openai.NewClientWithConfig(config)
	return client
}

func encodeImageToBase64(filePath string) string {
	// Open the image file
	file, _ := os.Open(filePath)

	defer file.Close()

	// Read the file contents into a byte slice
	fileInfo, _ := file.Stat()

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, _ = file.Read(buffer)

	// Encode the byte slice to a Base64 string
	base64String := base64.StdEncoding.EncodeToString(buffer)
	return "data:image/jpeg;base64," + base64String
}

func imageToText(link, data string) string {
	client := initAI()
	const prompt = `This is a screenshot of a webpage, I want you to transcribe this page's main content into text. 
Skip preamble and only return the result`

	imagePart := openai.ChatMessagePart{
		Type: openai.ChatMessagePartTypeImageURL,
		ImageURL: &openai.ChatMessageImageURL{
			URL: data,
		},
	}

	textPart := openai.ChatMessagePart{
		Type: openai.ChatMessagePartTypeText,
		Text: prompt,
	}
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4o,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:         openai.ChatMessageRoleUser,
					MultiContent: []openai.ChatMessagePart{imagePart, textPart},
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "ocr error"
	}

	var content = resp.Choices[0].Message.Content
	return fmt.Sprintf("Following text is a transcribe of %s:\n\n%s\n", link, content)
}

func askOpenAI(userID int64) string {
	client := initAI()
	chats := getChats(userID)

	var messages = []openai.ChatCompletionMessage{}
	for _, chat := range chats {
		m := openai.ChatCompletionMessage{
			Role:    chat.Role,
			Content: chat.Text,
		}
		messages = append(messages, m)
	}

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT4o,
			Messages: messages,
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return "chatgpt error"
	}

	var content = resp.Choices[0].Message.Content
	return content
}
