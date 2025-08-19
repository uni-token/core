package server

import (
	"bytes"
	"io"
	"net/http"
	"strings"
	"time"

	"uni-token-service/store"

	"github.com/gin-gonic/gin"
)

// SetupSiliconFlowAPI sets up SiliconFlow API endpoints
func SetupSiliconFlowAPI(router gin.IRouter) {
	api := router.Group("/siliconflow").Use(RequireUserLogin())
	{
		api.GET("/status", handleGetStatus)
		api.POST("/sms", handleSendSMS)
		api.POST("/login", handleSiliconLogin)
		api.POST("/logout", handleLogout)
		api.POST("/apikey/create", handleCreateAPIKey)
		api.POST("/payment/create", handleCreatePayment)
		api.GET("/payment/status", handleCheckPaymentStatus)
		api.GET("/auth/info", handleGetAuthInfo)
		api.POST("/auth/save", handleSaveAuth)
	}
}

// setCommonHeaders sets common HTTP headers for SiliconFlow requests
func setCommonHeaders(req *http.Request) {
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("Priority", "u=1, i")
	req.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"Linux"`)
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")
	req.Header.Set("Origin", "https://cloud.siliconflow.cn")
}

// handleSendSMS handles SMS sending request
func handleSendSMS(c *gin.Context) {
	// Read the raw request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", "https://account.siliconflow.cn/api/open/sms", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	setCommonHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send SMS request"})
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

// handleSiliconLogin handles user login request
func handleSiliconLogin(c *gin.Context) {
	// Read the raw request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", "https://account.siliconflow.cn/api/open/login/user", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	setCommonHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "text/plain;charset=UTF-8")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send login request"})
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Store cookies if login was successful
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		cookies := resp.Header["Set-Cookie"]
		var subjectID string

		if len(cookies) > 0 {
			// Parse Set-Cookie headers and extract only name=value pairs
			var cookiePairs []string
			for _, cookie := range cookies {
				// Split by semicolon and take only the first part (name=value)
				parts := strings.Split(cookie, ";")
				if len(parts) > 0 {
					nameValue := strings.TrimSpace(parts[0])
					if nameValue != "" {
						cookiePairs = append(cookiePairs, nameValue)
					}
				}
			}
			cookieStr := strings.Join(cookiePairs, "; ")

			// Make a request to /me to get the subject ID
			meReq, err := http.NewRequest("GET", "https://cloud.siliconflow.cn/me", nil)
			if err == nil {
				// Set headers for /me request
				meReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
				meReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
				meReq.Header.Set("Priority", "u=0, i")
				meReq.Header.Set("Referer", "https://account.siliconflow.cn/")
				meReq.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`)
				meReq.Header.Set("Sec-CH-UA-Mobile", "?0")
				meReq.Header.Set("Sec-CH-UA-Platform", `"Linux"`)
				meReq.Header.Set("Sec-Fetch-Dest", "document")
				meReq.Header.Set("Sec-Fetch-Mode", "navigate")
				meReq.Header.Set("Sec-Fetch-Site", "same-site")
				meReq.Header.Set("Sec-Fetch-User", "?1")
				meReq.Header.Set("Upgrade-Insecure-Requests", "1")
				meReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")
				meReq.Header.Set("Cookie", cookieStr)

				meResp, err := client.Do(meReq)
				if err == nil {
					defer meResp.Body.Close()
					// Extract X-Subject-ID from response headers
					if xSubjectID := meResp.Header.Get("X-Subject-ID"); xSubjectID != "" {
						subjectID = xSubjectID
					}
				}
			}

			session := store.SiliconFlowSession{
				Cookie:    cookieStr,
				SubjectID: subjectID,
				CreatedAt: time.Now(),
			}
			// Use a simple key for storage
			store.SiliconFlowSessions.Put("latest", session)
		}
	}

	c.Data(resp.StatusCode, "application/json", respBody)
}

// handleCreateAPIKey handles API key creation request
func handleCreateAPIKey(c *gin.Context) {
	// Retrieve the latest stored session
	session, err := store.SiliconFlowSessions.Get("latest")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session found. Please login first."})
		return
	}

	// Read the raw request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", "https://cloud.siliconflow.cn/biz-server/api/v1/apikey/create", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	setCommonHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Referer", "https://cloud.siliconflow.cn/me/account/ak")

	// Add X-Subject-ID if available
	if session.SubjectID != "" {
		httpReq.Header.Set("X-Subject-ID", session.SubjectID)
	}

	// Set the stored cookie
	httpReq.Header.Set("Cookie", session.Cookie)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send API key creation request"})
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

// handleCreatePayment handles payment QR code creation request
func handleCreatePayment(c *gin.Context) {
	// Retrieve the latest stored session
	session, err := store.SiliconFlowSessions.Get("latest")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session found. Please login first."})
		return
	}

	// Read the raw request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", "https://cloud.siliconflow.cn/biz-server/api/v1/pay/transactions", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	setCommonHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Referer", "https://cloud.siliconflow.cn/me/expensebill")

	// Add X-Subject-ID if available
	if session.SubjectID != "" {
		httpReq.Header.Set("X-Subject-ID", session.SubjectID)
	}

	// Set the stored cookie
	httpReq.Header.Set("Cookie", session.Cookie)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send payment creation request"})
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

// handleCheckPaymentStatus handles payment status checking request
func handleCheckPaymentStatus(c *gin.Context) {
	// Retrieve the latest stored session
	session, err := store.SiliconFlowSessions.Get("latest")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session found. Please login first."})
		return
	}

	// Get order parameter from query
	order := c.Query("order")
	if order == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order parameter is required"})
		return
	}

	// Create HTTP request
	url := "https://cloud.siliconflow.cn/biz-server/api/v1/pay/status?order=" + order
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	setCommonHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Referer", "https://cloud.siliconflow.cn/me/expensebill")

	// Add X-Subject-ID if available
	if session.SubjectID != "" {
		httpReq.Header.Set("X-Subject-ID", session.SubjectID)
	}

	// Set the stored cookie
	httpReq.Header.Set("Cookie", session.Cookie)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send payment status request"})
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

// handleGetStatus handles getting current login status and user info
func handleGetStatus(c *gin.Context) {
	// Retrieve the latest stored session
	session, err := store.SiliconFlowSessions.Get("latest")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code":    40001,
			"message": "No session found",
			"status":  false,
			"data":    nil,
		})
		return
	}

	// Make a request to /me to get user info
	client := &http.Client{}
	meReq, err := http.NewRequest("GET", "https://cloud.siliconflow.cn/me", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers for /me request
	meReq.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	meReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	meReq.Header.Set("Priority", "u=0, i")
	meReq.Header.Set("Referer", "https://account.siliconflow.cn/")
	meReq.Header.Set("Sec-CH-UA", `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`)
	meReq.Header.Set("Sec-CH-UA-Mobile", "?0")
	meReq.Header.Set("Sec-CH-UA-Platform", `"Linux"`)
	meReq.Header.Set("Sec-Fetch-Dest", "document")
	meReq.Header.Set("Sec-Fetch-Mode", "navigate")
	meReq.Header.Set("Sec-Fetch-Site", "same-site")
	meReq.Header.Set("Sec-Fetch-User", "?1")
	meReq.Header.Set("Upgrade-Insecure-Requests", "1")
	meReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")
	meReq.Header.Set("Cookie", session.Cookie)

	meResp, err := client.Do(meReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check login status"})
		return
	}
	defer meResp.Body.Close()

	// If the request is successful, user is logged in
	if meResp.StatusCode >= 200 && meResp.StatusCode < 300 {
		// Try to get user info from API
		userInfoReq, err := http.NewRequest("GET", "https://cloud.siliconflow.cn/biz-server/api/v1/user/info", nil)
		if err == nil {
			setCommonHeaders(userInfoReq)
			userInfoReq.Header.Set("Origin", "https://cloud.siliconflow.cn")
			userInfoReq.Header.Set("Referer", "https://cloud.siliconflow.cn/me/account/info")
			userInfoReq.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")
			if session.SubjectID != "" {
				userInfoReq.Header.Set("X-Subject-ID", session.SubjectID)
			}
			userInfoReq.Header.Set("Cookie", session.Cookie)

			userInfoResp, err := client.Do(userInfoReq)
			if err == nil {
				defer userInfoResp.Body.Close()
				if userInfoResp.StatusCode >= 200 && userInfoResp.StatusCode < 300 {
					userInfoBody, err := io.ReadAll(userInfoResp.Body)
					if err == nil {
						// Return the user info response directly since it already has the right format
						c.Data(http.StatusOK, "application/json", userInfoBody)
						return
					}
				}
			}
		}

		// If we can't get detailed user info, return basic login status with compatible format
		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"message": "User is logged in",
			"status":  true,
			"data": gin.H{
				"name": "User",
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    40001,
			"message": "Session expired or invalid",
			"status":  false,
			"data":    nil,
		})
	}
}

// handleLogout handles user logout request
func handleLogout(c *gin.Context) {
	// Remove the stored session
	err := store.SiliconFlowSessions.Delete("latest")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// handleGetAuthInfo handles getting real name authentication info
func handleGetAuthInfo(c *gin.Context) {
	// Retrieve the latest stored session
	session, err := store.SiliconFlowSessions.Get("latest")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session found. Please login first."})
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("GET", "https://cloud.siliconflow.cn/biz-server/api/v1/subject/auth/info", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	setCommonHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Referer", "https://cloud.siliconflow.cn/me/account/authentication/personal")

	// Add X-Subject-ID if available
	if session.SubjectID != "" {
		httpReq.Header.Set("X-Subject-ID", session.SubjectID)
	}

	// Set the stored cookie
	httpReq.Header.Set("Cookie", session.Cookie)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send auth info request"})
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

// handleSaveAuth handles saving real name authentication
func handleSaveAuth(c *gin.Context) {
	// Retrieve the latest stored session
	session, err := store.SiliconFlowSessions.Get("latest")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No session found. Please login first."})
		return
	}

	// Read the raw request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", "https://cloud.siliconflow.cn/biz-server/api/v1/subject/auth/save", bytes.NewBuffer(body))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers
	setCommonHeaders(httpReq)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Referer", "https://cloud.siliconflow.cn/me/account/authentication/personal")

	// Add X-Subject-ID if available
	if session.SubjectID != "" {
		httpReq.Header.Set("X-Subject-ID", session.SubjectID)
	}

	// Set the stored cookie
	httpReq.Header.Set("Cookie", session.Cookie)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send auth save request"})
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
