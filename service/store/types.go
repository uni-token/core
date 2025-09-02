package store

import "time"

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LLMKey struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`     // "manual", "siliconflow", etc.
	Protocol string `json:"protocol"` // "openai", "anthropic", etc.
	BaseURL  string `json:"baseUrl"`
	Token    string `json:"token"`
}

type AppPreset struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Keys      []string  `json:"keys"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type SiliconFlowSession struct {
	Cookie    string    `json:"cookie"`
	SubjectID string    `json:"subjectId"`
	CreatedAt time.Time `json:"createdAt"`
}

type AppInfo struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Key          string    `json:"key"`
	Granted      bool      `json:"granted"`
	CreatedAt    time.Time `json:"createdAt"`
	LastActiveAt time.Time `json:"lastActiveAt"`
}
