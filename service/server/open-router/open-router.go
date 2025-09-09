package openrouter

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"uni-token-service/store"

	"github.com/gin-gonic/gin"
)

type session struct {
	Key    string `json:"key"`
	UserId string `json:"userId"`
}

func loadSession() (session, error) {
	var s session
	data, err := store.Providers.Get("openRouter")
	if err != nil {
		return session{}, err
	}
	json.Unmarshal(data, &s)
	return s, nil
}

func saveSession(s session) error {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return err
	}
	store.Providers.Put("openRouter", jsonData)
	return nil
}

func SetupAPI(routes gin.IRoutes) {
	routes.POST("/authed", handleAuthed)
	routes.GET("/status", handleGetStatus)
	routes.GET("/key", handleGetKey)
	routes.POST("/logout", handleLogout)
}

func handleAuthed(c *gin.Context) {
	var req struct {
		Key    string `json:"key"`
		UserId string `json:"userId"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := saveSession(session{
		Key:    req.Key,
		UserId: req.UserId,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Status(http.StatusOK)
}

func handleGetStatus(c *gin.Context) {
	s, err := loadSession()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"authed": false,
		})
		return
	}

	req, err := http.NewRequest("GET", "https://openrouter.ai/api/v1/credits", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Authorization", "Bearer "+s.Key)
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch credits"})
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch credits"})
		return
	}
	var creditsResp struct {
		Data struct {
			TotalCredits float64 `json:"total_credits"`
			TotalUsage   float64 `json:"total_usage"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &creditsResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse credits"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"userId":  s.UserId,
		"credits": creditsResp.Data.TotalCredits,
		"usage":   creditsResp.Data.TotalUsage,
	})
}

func handleGetKey(c *gin.Context) {
	s, err := loadSession()
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	c.JSON(http.StatusOK, gin.H{"key": s.Key})
}

// handleLogout handles user logout request
func handleLogout(c *gin.Context) {
	// Remove the stored session
	err := store.Providers.Delete("openRouter")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}
	c.Status(http.StatusOK)
}
