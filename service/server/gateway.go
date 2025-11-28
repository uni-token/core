package server

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	"uni-token-service/logic"
	"uni-token-service/store"

	"github.com/gin-gonic/gin"
)

func SetupGatewayAPI(router *gin.Engine) {
	router.Any("/openai/*path", handleOpenAIProxy)
}

func handleOpenAIProxy(c *gin.Context) {
	path := c.Param("path")
	appId := ensureToken(c)

	appInfo, err := store.Apps.Get(appId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get app info"})
		return
	}

	if !appInfo.Granted {
		c.JSON(http.StatusForbidden, gin.H{"error": "App access not granted"})
		return
	}

	key, err := store.LLMKeys.Get(appInfo.Key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve key"})
		return
	}

	// Build target URL
	targetURL, err := url.JoinPath(key.BaseURL, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to build target URL"})
		return
	}

	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}

	// Read request body for usage tracking
	var requestBody []byte
	if c.Request.Body != nil {
		requestBody, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}

	// Extract model from request for usage tracking
	model := logic.ExtractModelFromRequest(requestBody)

	// Create new request
	req, err := http.NewRequest(c.Request.Method, targetURL, bytes.NewBuffer(requestBody))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy all headers except Authorization
	for key, values := range c.Request.Header {
		if strings.ToLower(key) != "authorization" {
			for _, value := range values {
				req.Header.Add(key, value)
			}
		}
	}

	// Set authorization header with key token
	req.Header.Set("Authorization", "Bearer "+key.Token)

	// Create HTTP client with no timeout for streaming
	client := &http.Client{
		Timeout: 0, // No timeout for streaming
	}

	resp, err := client.Do(req)
	if err != nil {
		// Record failed request
		logic.RecordUsage(appId, appInfo.Name, key.Name, model, path, 0, 0, 0, "error")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to proxy request"})
		return
	}
	defer resp.Body.Close()

	// Check if response is streaming (Server-Sent Events)
	contentType := resp.Header.Get("Content-Type")
	isStreaming := strings.Contains(contentType, "text/event-stream") ||
		(strings.Contains(contentType, "text/plain") && resp.Header.Get("Transfer-Encoding") == "chunked")

	// Copy all response headers
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Set status code
	c.Status(resp.StatusCode)

	if isStreaming {
		// Handle streaming response
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Writer.Flush()

		// Create usage extractor for streaming
		usageExtractor := logic.NewStreamingUsageExtractor(model)
		usageExtractor.SetContext(appId, appInfo.Name, key.Name, path)

		// Stream response body
		buffer := make([]byte, 4096)
		for {
			n, err := resp.Body.Read(buffer)
			if n > 0 {
				// Extract usage from streaming chunks
				usageExtractor.ProcessChunk(buffer[:n])

				if _, writeErr := c.Writer.Write(buffer[:n]); writeErr != nil {
					break
				}
				c.Writer.Flush()
			}
			if err != nil {
				if err != io.EOF {
					// Log error but don't return JSON as we're already streaming
				}
				break
			}
		}

		// Record streaming usage
		status := "success"
		if resp.StatusCode >= 400 {
			status = "error"
		}

		usageExtractor.RecordUsage(status)
	} else {
		// Handle non-streaming response
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			// Record failed request
			logic.RecordUsage(appId, appInfo.Name, key.Name, model, path, 0, 0, 0, "error")
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
			return
		}

		// Extract usage from response
		usageData := logic.ExtractUsageFromResponse(responseBody)

		// Record usage
		status := "success"
		if resp.StatusCode >= 400 {
			status = "error"
		}

		// Use model from response if available, otherwise use request model
		finalModel := model
		if usageData.Model != "unknown" {
			finalModel = usageData.Model
		}

		logic.RecordUsage(appId, appInfo.Name, key.Name, finalModel, path,
			usageData.PromptTokens, usageData.OutputTokens, usageData.Cost, status)

		// Write response body
		c.Writer.Write(responseBody)
	}
}

func ensureToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Authorization token is required"})
		c.Abort()
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}
