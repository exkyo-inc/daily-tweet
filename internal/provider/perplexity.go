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

// PerplexityProvider provides anniversaries using Perplexity API
type PerplexityProvider struct {
	apiKey string
	client *http.Client
}

// NewPerplexityProvider creates a new Perplexity provider
func NewPerplexityProvider(apiKey string) *PerplexityProvider {
	return &PerplexityProvider{
		apiKey: apiKey,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Name returns the provider name
func (p *PerplexityProvider) Name() string {
	return "Perplexity"
}

// PerplexityRequest represents the request structure for Perplexity API
type PerplexityRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// PerplexityResponse represents the response from Perplexity API
type PerplexityResponse struct {
	Choices []Choice `json:"choices"`
}

// GetAnniversaries returns anniversaries for the given date using Perplexity
func (p *PerplexityProvider) GetAnniversaries(date time.Time) ([]model.Anniversary, error) {
	prompt := fmt.Sprintf("今日は%d年%d月%d日です。この日は何の日ですか？歴史的な出来事や記念日を教えてください。簡潔に答えてください。",
		date.Year(), date.Month(), date.Day())

	request := PerplexityRequest{
		Model: "llama-3.1-sonar-small-128k-online",
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

	req, err := http.NewRequest("POST", "https://api.perplexity.ai/chat/completions", bytes.NewBuffer(jsonData))
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

	var response PerplexityResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(response.Choices) == 0 {
		return nil, fmt.Errorf("no response from Perplexity")
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
