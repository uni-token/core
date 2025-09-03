package deepSeek

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"uni-token-service/store"

	"github.com/gin-gonic/gin"
)

type session struct {
	Token string `json:"token"`
}

func loadSession() (session, error) {
	var s session
	data, err := store.Providers.Get("deepSeek")
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
	store.Providers.Put("deepSeek", jsonData)
	return nil
}

func SetupAPI(routes gin.IRoutes) {
	// Login
	routes.POST("/login", handleLogin)

	// Send SMS
	routes.POST("/sms", func(ctx *gin.Context) {
		forward(ctx, false, "https://platform.deepseek.com/auth-api/v0/users/create_sms_verification_code")
	})

	// User status
	routes.GET("/status", func(ctx *gin.Context) {
		_, err := loadSession()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{})
			return
		}

		forward(ctx, true, "https://platform.deepseek.com/auth-api/v0/users/current")
	})

	// Logout
	routes.POST("/logout", func(ctx *gin.Context) {
		// Clear stored session
		store.Providers.Delete("deepSeek")
		ctx.JSON(http.StatusOK, gin.H{"success": true})
	})

	// Identity verification
	routes.GET("/auth/info", func(ctx *gin.Context) {
		forward(ctx, true, "https://platform.deepseek.com/api/v1/my_identity_verification")
	})

	routes.POST("/auth/save", func(ctx *gin.Context) {
		forward(ctx, true, "https://platform.deepseek.com/api/v1/identity_verify")
	})

	// Payment
	routes.POST("/payment/create", func(ctx *gin.Context) {
		forward(ctx, true, "https://platform.deepseek.com/api/v1/payments")
	})

	routes.POST("/payment/status", func(ctx *gin.Context) {
		orderId := ctx.Query("order")
		if orderId == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "order parameter is required"})
			return
		}
		target := "https://platform.deepseek.com/api/v1/payments/" + orderId + "/capture"
		forward(ctx, true, target)
	})

	// API Key management
	routes.POST("/apikey/create", func(ctx *gin.Context) {
		forward(ctx, true, "https://platform.deepseek.com/api/v0/users/edit_api_keys")
	})
}

// setCommonHeaders sets common HTTP headers for SiliconFlow requests
func setCommonHeaders(req *http.Request) {
	req.Header.Set("accept", "*/*")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://platform.deepseek.com/top_up")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")
}

func forward(c *gin.Context, requireLogin bool, target string) {
	method := c.Request.Method

	// Retrieve the latest stored session
	var session session
	var err error
	if requireLogin {
		session, err = loadSession()
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No session found. Please login first."})
			return
		}
	}

	// Read request body for POST requests
	var reqBody io.Reader
	if method == "POST" {
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
			return
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request with the specified method
	httpReq, err := http.NewRequest(method, target, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	setCommonHeaders(httpReq)
	if requireLogin {
		httpReq.Header.Set("authorization", "Bearer "+session.Token)
	}

	// Copy Content-Type from original request if it's a POST
	if method == "POST" {
		if contentType := c.GetHeader("Content-Type"); contentType != "" {
			httpReq.Header.Set("Content-Type", contentType)
		} else {
			httpReq.Header.Set("Content-Type", "application/json")
		}
	}

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(resp.StatusCode, "application/json", respBody)
}

func handleLogin(c *gin.Context) {
	var req struct {
		Phone    string `json:"phone"`
		Code     string `json:"code"`
		AreaCode string `json:"area_code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Prepare form data for login
	formData := url.Values{}
	formData.Set("mobile_number", req.Phone)
	formData.Set("sms_code", req.Code)
	formData.Set("area_code", req.AreaCode)

	// Create login request
	httpReq, err := http.NewRequest("POST", "https://platform.deepseek.com/auth-api/v0/users/login_by_mobile_sms", strings.NewReader(formData.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create login request"})
		return
	}

	// Set headers for login
	setCommonHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Make login request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login"})
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read login response"})
		return
	}

	// Parse response to extract token
	var loginResp map[string]interface{}
	if err := json.Unmarshal(respBody, &loginResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse login response"})
		return
	}

	// Extract token and save session
	if data, ok := loginResp["data"].(map[string]interface{}); ok {
		if token, ok := data["token"].(string); ok {
			session := session{Token: token}
			if err := saveSession(session); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
				return
			}
		}
	}

	c.Data(resp.StatusCode, "application/json", respBody)
}
