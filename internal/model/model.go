package model

import "time"

// Anniversary represents a historical event or commemoration
type Anniversary struct {
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
}

// Config holds the application configuration
type Config struct {
	DiscordWebhookURL string `json:"discord_webhook_url"`
	OpenAIAPIKey      string `json:"openai_api_key"`
	PerplexityAPIKey  string `json:"perplexity_api_key"`
	DryRun            bool   `json:"dry_run"`
}

// Provider interface for different anniversary sources
type Provider interface {
	GetAnniversaries(date time.Time) ([]Anniversary, error)
	Name() string
}
