package provider

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cocacola/daily-tweet/internal/model"
)

// StaticProvider provides anniversaries from a static CSV file
type StaticProvider struct {
	filePath string
}

// NewStaticProvider creates a new static provider
func NewStaticProvider(filePath string) *StaticProvider {
	return &StaticProvider{
		filePath: filePath,
	}
}

// Name returns the provider name
func (p *StaticProvider) Name() string {
	return "Static CSV"
}

// GetAnniversaries returns anniversaries for the given date from CSV file
func (p *StaticProvider) GetAnniversaries(date time.Time) ([]model.Anniversary, error) {
	file, err := os.Open(p.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", p.filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // CSV format
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file: %w", err)
	}

	var anniversaries []model.Anniversary
	targetDate := fmt.Sprintf("%d-%02d", int(date.Month()), date.Day())

	// Skip header row
	for i, record := range records {
		if i == 0 || len(record) < 6 {
			continue
		}

		recordDate := strings.TrimSpace(record[0])
		title := strings.TrimSpace(record[1])
		description := strings.TrimSpace(record[5]) // 記念日の説明は6列目

		if recordDate == targetDate {
			anniversary := model.Anniversary{
				Date:        date,
				Title:       title,
				Description: description,
				Source:      p.Name(),
			}
			anniversaries = append(anniversaries, anniversary)
		}
	}

	return anniversaries, nil
}
