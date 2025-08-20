package discovery

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

type uniTokenDetectionResponse struct {
	UniToken bool `json:"__uni_token"`
}

func IsServiceRunning() bool {
	filePath := getServiceJsonPath()

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}

	var info ServiceInfo
	if err := json.Unmarshal(fileContent, &info); err != nil {
		return false
	}

	if info.URL == "" {
		return false
	}

	// Verify the service is actually running
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(info.URL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var detection uniTokenDetectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&detection); err != nil {
		return false
	}

	return detection.UniToken
}
