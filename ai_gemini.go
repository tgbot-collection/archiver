// archiver - ai.go
// 2024-02-20 20:58
// Benny <benny.think@gmail.com>

package main

import (
	"bytes"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type Response struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content ContentData `json:"content"`
}

type ContentData struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Role  string `json:"role"`
	Parts []Part `json:"parts"`
}

type Data struct {
	Contents []Content `json:"contents"`
}

func askGemini(userID int64) string {
	var data Data
	chats := getChats(userID)
	if len(chats) > 0 {
		for _, chat := range chats {
			part := Part{Text: chat.Text}
			content := Content{
				Role:  chat.Role,
				Parts: []Part{part},
			}
			data.Contents = append(data.Contents, content)
		}
	}

	jsonData, _ := json.Marshal(data)
	log.Infoln(string(jsonData))
	resp, err := http.Post(geminiURL, "application/json", bytes.NewBuffer(jsonData))
	defer resp.Body.Close()
	if err != nil || resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		text := string(body)
		log.Errorf("Request failed, %s", text)
		return text
	}

	var response Response
	_ = json.NewDecoder(resp.Body).Decode(&response)
	return response.Candidates[0].Content.Parts[0].Text
}
