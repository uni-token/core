package server

import (
	"bytes"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupProxyAPI(router gin.IRouter) {
	router.POST("/proxy", handleProxy, RequireUserLogin())
	router.GET("/proxy/:base/*paths", handleSimpleProxy)
}

type ProxyRequest struct {
	Method  string            `json:"method" binding:"required"`
	Url     string            `json:"url" binding:"required"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

type ProxyResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func handleProxy(c *gin.Context) {
	var req ProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{}
	proxyReq, err := http.NewRequest(req.Method, req.Url, bytes.NewReader([]byte(req.Body)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request: " + err.Error()})
		return
	}

	for key, value := range req.Headers {
		proxyReq.Header.Set(key, value)
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform proxy request: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	var respBody bytes.Buffer
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read proxy response body: " + err.Error()})
		return
	}

	respHeaders := make(map[string]string)
	for key, values := range resp.Header {
		if len(values) > 0 {
			respHeaders[key] = values[0]
		}
	}

	c.JSON(http.StatusOK, ProxyResponse{
		Status:  resp.StatusCode,
		Headers: respHeaders,
		Body:    respBody.String(),
	})
}

func handleSimpleProxy(c *gin.Context) {
	base, err := base64.StdEncoding.DecodeString(c.Param("base"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid base64 URL"})
		return
	}

	paths := c.Param("paths")
	targetUrl := string(base) + paths

	client := &http.Client{}
	proxyReq, err := http.NewRequest(c.Request.Method, targetUrl, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request: " + err.Error()})
		return
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform proxy request: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	var respBody bytes.Buffer
	_, err = respBody.ReadFrom(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read proxy response body: " + err.Error()})
		return
	}

	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody.Bytes())
}
