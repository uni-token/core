package logic

import (
	"encoding/json"
	"strings"
	"uni-token-service/store"
)

// ModelPricing represents the pricing structure for a model
type ModelPricing struct {
	PromptRate float64 // USD per 1K tokens
	OutputRate float64 // USD per 1K tokens
}

// GetModelPricing returns pricing information for a given model
func GetModelPricing(model string) ModelPricing {
	modelLower := strings.ToLower(model)

	switch {
	case strings.Contains(modelLower, "gpt-4o"):
		return ModelPricing{
			PromptRate: 0.005, // $0.005 per 1K prompt tokens
			OutputRate: 0.015, // $0.015 per 1K output tokens
		}
	case strings.Contains(modelLower, "gpt-4"):
		return ModelPricing{
			PromptRate: 0.03, // $0.03 per 1K prompt tokens
			OutputRate: 0.06, // $0.06 per 1K output tokens
		}
	case strings.Contains(modelLower, "gpt-3.5"):
		return ModelPricing{
			PromptRate: 0.0015, // $0.0015 per 1K prompt tokens
			OutputRate: 0.002,  // $0.002 per 1K output tokens
		}
	case strings.Contains(modelLower, "claude"):
		return ModelPricing{
			PromptRate: 0.008, // $0.008 per 1K prompt tokens
			OutputRate: 0.024, // $0.024 per 1K output tokens
		}
	default:
		return ModelPricing{
			PromptRate: 0.001, // Default rate
			OutputRate: 0.002,
		}
	}
}

// CalculateCost calculates the cost based on model and token usage
func CalculateCost(model string, promptTokens, outputTokens int) float64 {
	pricing := GetModelPricing(model)

	promptCost := float64(promptTokens) / 1000.0 * pricing.PromptRate
	outputCost := float64(outputTokens) / 1000.0 * pricing.OutputRate

	return promptCost + outputCost
}

// ExtractModelFromRequest extracts model name from request body
func ExtractModelFromRequest(requestBody []byte) string {
	if len(requestBody) == 0 {
		return "unknown"
	}

	var req map[string]interface{}
	if err := json.Unmarshal(requestBody, &req); err != nil {
		return "unknown"
	}

	if model, ok := req["model"].(string); ok {
		return model
	}

	return "unknown"
}

// UsageData represents token usage information
type UsageData struct {
	PromptTokens int
	OutputTokens int
	Cost         float64
	Model        string
}

// ExtractUsageFromResponse extracts token usage from API response
func ExtractUsageFromResponse(responseBody []byte) UsageData {
	var resp map[string]interface{}
	if err := json.Unmarshal(responseBody, &resp); err != nil {
		return UsageData{}
	}

	usage, ok := resp["usage"].(map[string]interface{})
	if !ok {
		return UsageData{}
	}

	var promptTokens, outputTokens int
	if prompt, ok := usage["prompt_tokens"].(float64); ok {
		promptTokens = int(prompt)
	}

	if completion, ok := usage["completion_tokens"].(float64); ok {
		outputTokens = int(completion)
	}

	// Get model from response or use default
	model := "unknown"
	if modelStr, ok := resp["model"].(string); ok {
		model = modelStr
	}

	cost := CalculateCost(model, promptTokens, outputTokens)

	return UsageData{
		PromptTokens: promptTokens,
		OutputTokens: outputTokens,
		Cost:         cost,
		Model:        model,
	}
}

// RecordUsage records token usage to the store
func RecordUsage(appID, appName, key, model, endpoint string, promptTokens, outputTokens int, cost float64, status string) error {
	return store.RecordUsage(appID, appName, key, model, endpoint, promptTokens, outputTokens, cost, status)
}

// StreamingUsageExtractor extracts usage data from streaming responses
type StreamingUsageExtractor struct {
	AppID        string
	AppName      string
	Key          string
	Model        string
	Endpoint     string
	PromptTokens int
	OutputTokens int
	TotalTokens  int
	buffer       string
}

// NewStreamingUsageExtractor creates a new streaming usage extractor
func NewStreamingUsageExtractor(model string) *StreamingUsageExtractor {
	return &StreamingUsageExtractor{
		Model: model,
	}
}

// ProcessChunk processes a chunk of streaming data to extract usage information
func (s *StreamingUsageExtractor) ProcessChunk(chunk []byte) {
	s.buffer += string(chunk)

	// Look for usage information in SSE format
	lines := strings.Split(s.buffer, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "data: ") {
			dataStr := strings.TrimPrefix(line, "data: ")
			if dataStr == "[DONE]" {
				continue
			}

			var data map[string]interface{}
			if err := json.Unmarshal([]byte(dataStr), &data); err == nil {
				// Extract usage from streaming chunk
				if usage, ok := data["usage"].(map[string]interface{}); ok {
					if promptTokens, ok := usage["prompt_tokens"].(float64); ok {
						s.PromptTokens = int(promptTokens)
					}
					if completionTokens, ok := usage["completion_tokens"].(float64); ok {
						s.OutputTokens = int(completionTokens)
					}
					if totalTokens, ok := usage["total_tokens"].(float64); ok {
						s.TotalTokens = int(totalTokens)
					}
				}

				// Update model if available in streaming response
				if model, ok := data["model"].(string); ok && model != "" {
					s.Model = model
				}
			}
		}
	}

	// Keep last incomplete line in buffer
	if len(lines) > 0 {
		s.buffer = lines[len(lines)-1]
	}
}

// RecordUsage records the collected usage data for streaming
func (s *StreamingUsageExtractor) RecordUsage(status string) error {
	cost := CalculateCost(s.Model, s.PromptTokens, s.OutputTokens)
	return RecordUsage(s.AppID, s.AppName, s.Key, s.Model, s.Endpoint, s.PromptTokens, s.OutputTokens, cost, status)
}

// SetContext sets the context information for the streaming extractor
func (s *StreamingUsageExtractor) SetContext(appID, appName, key, endpoint string) {
	s.AppID = appID
	s.AppName = appName
	s.Key = key
	s.Endpoint = endpoint
}

// GetUsageData returns the collected usage data
func (s *StreamingUsageExtractor) GetUsageData() UsageData {
	cost := CalculateCost(s.Model, s.PromptTokens, s.OutputTokens)
	return UsageData{
		PromptTokens: s.PromptTokens,
		OutputTokens: s.OutputTokens,
		Cost:         cost,
		Model:        s.Model,
	}
}
