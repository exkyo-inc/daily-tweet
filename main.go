package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cocacola/daily-tweet/internal/model"
	"github.com/cocacola/daily-tweet/internal/provider"
)

func postToDiscord(webhookURL string, message string, dryRun bool) error {
	if dryRun || webhookURL == "" {
		fmt.Println("[DryRun]", message)
		return nil
	}

	payload := struct {
		Content string `json:"content"`
	}{Content: message}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to post to webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook returned status %s", resp.Status)
	}

	return nil
}

func main() {
	fmt.Println("Today-in-History Bot")

	// Load configuration from environment variables
	config := &model.Config{
		DiscordWebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		OpenAIAPIKey:      os.Getenv("OPENAI_API_KEY"),
		PerplexityAPIKey:  os.Getenv("PERPLEXITY_API_KEY"),
		DryRun:            os.Getenv("DRY_RUN") == "true",
	}

	// Initialize providers
	staticProvider := provider.NewStaticProvider("data/anniversaries.csv")

	var providers []model.Provider
	providers = append(providers, staticProvider)

	if config.OpenAIAPIKey != "" {
		chatgptProvider := provider.NewChatGPTProvider(config.OpenAIAPIKey)
		providers = append(providers, chatgptProvider)
	}

	if config.PerplexityAPIKey != "" {
		perplexityProvider := provider.NewPerplexityProvider(config.PerplexityAPIKey)
		providers = append(providers, perplexityProvider)
	}

	// Generate posts for the next 7 days
	today := time.Now()
	for i := 0; i < 7; i++ {
		targetDate := today.AddDate(0, 0, i)

		fmt.Printf("\n=== %s ===\n", targetDate.Format("2006-01-02"))

		for _, p := range providers {
			anniversaries, err := p.GetAnniversaries(targetDate)
			if err != nil {
				log.Printf("Error getting anniversaries from %s: %v", p.Name(), err)
				continue
			}

			for _, anniversary := range anniversaries {
				msg := fmt.Sprintf("[%s] %s: %s", anniversary.Source, anniversary.Title, anniversary.Description)
				if err := postToDiscord(config.DiscordWebhookURL, msg, config.DryRun); err != nil {
					log.Printf("Error posting to Discord: %v", err)
				}
			}
		}
	}

	log.Println("Bot completed successfully")
}
