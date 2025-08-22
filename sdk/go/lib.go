package uniToken

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// UniTokenOptions represents the options for requesting a UniToken OpenAI token
type UniTokenOptions struct {
	// AppName is the name of the application requesting the OpenAI token
	AppName string `json:"appName"`
	// Description is a brief description of the application
	Description string `json:"description"`
	// SavedAPIKey is an optional saved API key, if the user has previously granted permission
	SavedAPIKey string `json:"savedApiKey,omitempty"`
}

// UniTokenResult represents the result of a UniToken request
type UniTokenResult struct {
	BaseURL string `json:"baseUrl"`
	// Empty if the user did not grant permission
	APIKey string `json:"apiKey"`
}

// serviceInfo represents the structure of service.json file
type serviceInfo struct {
	URL string `json:"url"`
}

type appRegisterRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UID         string `json:"uid,omitempty"`
}

type appRegisterResponse struct {
	Token string `json:"token"`
}

type uniTokenDetectionResponse struct {
	UniToken bool `json:"__uni_token"`
}

// setupServiceRootPath creates and returns the service root directory path
func setupServiceRootPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}

	var root string
	if runtime.GOOS == "windows" {
		root = filepath.Join(homeDir, "AppData", "Local", "UniToken")
	} else {
		root = filepath.Join(homeDir, ".local", "share", "uni-token")
	}

	if err := os.MkdirAll(root, 0755); err != nil {
		return "", fmt.Errorf("failed to create root directory: %w", err)
	}

	return root, nil
}

// detectRunningURLFromFile detects if the service is running by reading service.json
func detectRunningURLFromFile(rootPath string) (string, error) {
	filePath := filepath.Join(rootPath, "service.json")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", nil
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return "", nil
	}

	var info serviceInfo
	if err := json.Unmarshal(fileContent, &info); err != nil {
		return "", nil
	}

	if info.URL == "" {
		return "", nil
	}

	// Verify the service is actually running
	client := &http.Client{Timeout: 10 * time.Minute}
	resp, err := client.Get(info.URL)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	var detection uniTokenDetectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&detection); err != nil {
		return "", nil
	}

	if detection.UniToken {
		return info.URL, nil
	}

	return "", nil
}

// startService starts the UniToken service
func startService(rootPath string) (string, error) {
	var execPath string
	if runtime.GOOS == "windows" {
		execPath = filepath.Join(rootPath, "service.exe")
	} else {
		execPath = filepath.Join(rootPath, "service")
	}

	if _, err := os.Stat(execPath); os.IsNotExist(err) {
		if err := downloadService(execPath); err != nil {
			return "", fmt.Errorf("failed to download service: %w", err)
		}
	}

	cmd := exec.Command(execPath, "setup")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to start service: %w", err)
	}

	// Wait a moment for the service to start and write the URL file
	time.Sleep(2 * time.Second)

	serverURL, err := detectRunningURLFromFile(rootPath)
	if err != nil || serverURL == "" {
		return "", fmt.Errorf("service started but URL not detected")
	}

	return serverURL, nil
}

// downloadService downloads the appropriate service binary for the current platform
func downloadService(execPath string) error {
	platformMap := map[string]string{
		"linux":   "service-linux-amd64",
		"darwin":  "service-darwin-amd64",
		"windows": "service-windows-amd64.exe",
	}

	filename, exists := platformMap[runtime.GOOS]
	if !exists {
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	url := fmt.Sprintf("https://uni-token.app/release/%s", filename)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: HTTP %d - %s", resp.StatusCode, resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := os.WriteFile(execPath, data, 0755); err != nil {
		return fmt.Errorf("failed to write service binary: %w", err)
	}

	return nil
}

// RequestUniTokenOpenAI requests user for OpenAI token via UniToken service
// Returns the baseURL and apiKey. apiKey is empty if the user does not grant permission
func RequestUniTokenOpenAI(options UniTokenOptions) (UniTokenResult, error) {
	rootPath, err := setupServiceRootPath()
	if err != nil {
		return UniTokenResult{}, fmt.Errorf("failed to setup service root path: %w", err)
	}
	serverURL, err := detectRunningURLFromFile(rootPath)

	if err != nil || serverURL == "" {
		serverURL, err = startService(rootPath)
		if err != nil {
			return UniTokenResult{}, fmt.Errorf("failed to start service: %w", err)
		}
	}

	requestBody := appRegisterRequest{
		Name:        options.AppName,
		Description: options.Description,
		UID:         options.SavedAPIKey,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return UniTokenResult{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(
		fmt.Sprintf("%sapp/register", serverURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return UniTokenResult{}, fmt.Errorf("registration request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden {
		return UniTokenResult{
			BaseURL: fmt.Sprintf("%sopenai/", serverURL),
			APIKey:  "",
		}, nil // User denied permission
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return UniTokenResult{}, fmt.Errorf("registration failed: HTTP %d - %s", resp.StatusCode, string(body))
	}

	var registerResp appRegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
		return UniTokenResult{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return UniTokenResult{
		BaseURL: fmt.Sprintf("%sopenai/", serverURL),
		APIKey:  registerResp.Token,
	}, nil
}
