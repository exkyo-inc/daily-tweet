package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/cocacola/daily-tweet/internal/model"
)

// ChatGPTProvider provides anniversaries using OpenAI ChatGPT API
type ChatGPTProvider struct {
	apiKey string
	client *http.Client
}

// NewChatGPTProvider creates a new ChatGPT provider
func NewChatGPTProvider(apiKey string) *ChatGPTProvider {
	return &ChatGPTProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Name returns the provider name
func (p *ChatGPTProvider) Name() string {
	return "ChatGPT"
}

// ChatGPTRequest represents the request structure for OpenAI API
type ChatGPTRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatGPTResponse represents the response from OpenAI API
type ChatGPTResponse struct {
	Choices []Choice `json:"choices"`
}

// Choice represents a response choice
type Choice struct {
	Message Message `json:"message"`
}

// GetAnniversaries returns anniversaries for the given date using ChatGPT
func (p *ChatGPTProvider) GetAnniversaries(date time.Time) ([]model.Anniversary, error) {
	prompt := fmt.Sprintf("今日は%d年%d月%d日です。この日は何の日ですか？歴史的な出来事や記念日を教えてください。簡潔に答えてください。",
		date.Year(), date.Month(), date.Day())

	request := ChatGPTRequest{
		Model: "gpt-3.5-turbo",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response ChatGPTResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from ChatGPT")
	}

	content := response.Choices[0].Message.Content
	anniversary := model.Anniversary{
		Date:        date,
		Title:       fmt.Sprintf("%d月%d日の出来事", date.Month(), date.Day()),
		Description: content,
		Source:      p.Name(),
	}

	return []model.Anniversary{anniversary}, nil
}
