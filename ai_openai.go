package main

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func mapRole(originalRole string) string {
	switch originalRole {
	case "USER":
		return "user"
	case "MODEL":
		return "assistant"
	default:
		return "user"
	}
}

const openAI = "https://gptmos.com/v1/chat/completions"

func askOpenAI(userID int64) string {
	var chatReq ChatRequest
	chats := getChats(userID)
	if len(chats) > 0 {
		for _, chat := range chats {
			chatReq.Messages = append(chatReq.Messages, ChatMessage{
				Role:    mapRole(chat.Role),
				Content: chat.Text,
			})
		}
	}

	chatReq.Model = "gpt-4-0125-preview"

	jsonData, err := json.Marshal(chatReq)
	if err != nil {
		log.Errorf("Failed to marshal request: %v", err)
		return err.Error()
	}
	log.Infoln(string(jsonData))

	// Replace with the actual URL and add authorization headers as needed
	req, _ := http.NewRequest("POST", openAI, bytes.NewBuffer(jsonData))
	req.Header.Set("Authorization", "Bearer "+OpenAIKey)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Failed to make request: %v", err)
		return err.Error()
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		text := string(body)
		log.Errorf("Request failed, %s", text)
		return text
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		log.Errorf("Failed to decode response: %v", err)
		return err.Error()
	}

	if len(chatResp.Choices) > 0 && len(chatResp.Choices[0].Message.Content) > 0 {
		return chatResp.Choices[0].Message.Content
	}

	return "no response found"
}
