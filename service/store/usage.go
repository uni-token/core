package store

import (
	"time"
)

// TokenUsage represents a token usage record
type TokenUsage struct {
	ID           string    `json:"id"`
	AppID        string    `json:"appId"`
	AppName      string    `json:"appName"`
	Provider     string    `json:"provider"`
	Model        string    `json:"model"`
	PromptTokens int       `json:"promptTokens"`
	OutputTokens int       `json:"outputTokens"`
	TotalTokens  int       `json:"totalTokens"`
	Cost         float64   `json:"cost"`
	Timestamp    time.Time `json:"timestamp"`
	Endpoint     string    `json:"endpoint"`
	Status       string    `json:"status"`
}

// UsageStats represents aggregated usage statistics
type UsageStats struct {
	TotalTokens   int                      `json:"totalTokens"`
	TotalCost     float64                  `json:"totalCost"`
	TotalRequests int                      `json:"totalRequests"`
	ByApp         map[string]AppUsage      `json:"byApp"`
	ByProvider    map[string]ProviderUsage `json:"byProvider"`
	ByModel       map[string]ModelUsage    `json:"byModel"`
	RecentUsages  []TokenUsage             `json:"recentUsages"`
}

type AppUsage struct {
	AppName      string  `json:"appName"`
	TotalTokens  int     `json:"totalTokens"`
	TotalCost    float64 `json:"totalCost"`
	RequestCount int     `json:"requestCount"`
}

type ProviderUsage struct {
	TotalTokens  int     `json:"totalTokens"`
	TotalCost    float64 `json:"totalCost"`
	RequestCount int     `json:"requestCount"`
}

type ModelUsage struct {
	Provider     string  `json:"provider"`
	TotalTokens  int     `json:"totalTokens"`
	TotalCost    float64 `json:"totalCost"`
	RequestCount int     `json:"requestCount"`
}

// GetUsageStats calculates and returns usage statistics
func GetUsageStats(days int) (*UsageStats, error) {
	usages, err := Usage.List()
	if err != nil {
		return nil, err
	}

	cutoff := time.Now().AddDate(0, 0, -days)
	stats := &UsageStats{
		ByApp:        make(map[string]AppUsage),
		ByProvider:   make(map[string]ProviderUsage),
		ByModel:      make(map[string]ModelUsage),
		RecentUsages: make([]TokenUsage, 0, len(usages)),
	}

	for _, usage := range usages {
		if usage.Timestamp.Before(cutoff) {
			continue
		}

		stats.TotalTokens += usage.TotalTokens
		stats.TotalCost += usage.Cost
		stats.TotalRequests++

		// By App
		if appUsage, exists := stats.ByApp[usage.AppID]; exists {
			appUsage.TotalTokens += usage.TotalTokens
			appUsage.TotalCost += usage.Cost
			appUsage.RequestCount++
			stats.ByApp[usage.AppID] = appUsage
		} else {
			stats.ByApp[usage.AppID] = AppUsage{
				AppName:      usage.AppName,
				TotalTokens:  usage.TotalTokens,
				TotalCost:    usage.Cost,
				RequestCount: 1,
			}
		}

		// By Provider
		if providerUsage, exists := stats.ByProvider[usage.Provider]; exists {
			providerUsage.TotalTokens += usage.TotalTokens
			providerUsage.TotalCost += usage.Cost
			providerUsage.RequestCount++
			stats.ByProvider[usage.Provider] = providerUsage
		} else {
			stats.ByProvider[usage.Provider] = ProviderUsage{
				TotalTokens:  usage.TotalTokens,
				TotalCost:    usage.Cost,
				RequestCount: 1,
			}
		}

		// By Model
		modelKey := usage.Provider + "/" + usage.Model
		if modelUsage, exists := stats.ByModel[modelKey]; exists {
			modelUsage.TotalTokens += usage.TotalTokens
			modelUsage.TotalCost += usage.Cost
			modelUsage.RequestCount++
			stats.ByModel[modelKey] = modelUsage
		} else {
			stats.ByModel[modelKey] = ModelUsage{
				Provider:     usage.Provider,
				TotalTokens:  usage.TotalTokens,
				TotalCost:    usage.Cost,
				RequestCount: 1,
			}
		}

		stats.RecentUsages = append(stats.RecentUsages, usage)
	}

	// Sort recent usages by timestamp (most recent first)
	for i := 0; i < len(stats.RecentUsages)-1; i++ {
		for j := i + 1; j < len(stats.RecentUsages); j++ {
			if stats.RecentUsages[i].Timestamp.Before(stats.RecentUsages[j].Timestamp) {
				stats.RecentUsages[i], stats.RecentUsages[j] = stats.RecentUsages[j], stats.RecentUsages[i]
			}
		}
	}

	// Limit to last 100 records
	if len(stats.RecentUsages) > 100 {
		stats.RecentUsages = stats.RecentUsages[:100]
	}

	return stats, nil
}

// RecordUsage records a new token usage
func RecordUsage(appID, appName, provider, model, endpoint string, promptTokens, outputTokens int, cost float64, status string) error {
	usage := TokenUsage{
		ID:           generateID(),
		AppID:        appID,
		AppName:      appName,
		Provider:     provider,
		Model:        model,
		PromptTokens: promptTokens,
		OutputTokens: outputTokens,
		TotalTokens:  promptTokens + outputTokens,
		Cost:         cost,
		Timestamp:    time.Now(),
		Endpoint:     endpoint,
		Status:       status,
	}

	return Usage.Put(usage.ID, usage)
}

func generateID() string {
	return time.Now().Format("20060102150405") + generateRandomString(6)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(result)
}
