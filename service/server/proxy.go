package server

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupProxyAPI(router gin.IRouter) {
	api := router.Group("/proxy").Use(RequireUserLogin())
	{
		api.GET("", handleProxy)
		api.POST("", handleProxy)
		api.DELETE("", handleProxy)
		api.PUT("", handleProxy)
	}
}

type ProxyRequest struct {
	Url     string            `json:"url" binding:"required"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func handleProxy(c *gin.Context) {
	println("Handle proxy")

	var req ProxyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client := &http.Client{}
	proxyReq, err := http.NewRequest(c.Request.Method, req.Url, bytes.NewReader([]byte(req.Body)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
		return
	}

	for key, value := range req.Headers {
		proxyReq.Header.Set(key, value)
	}

	resp, err := client.Do(proxyReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to perform proxy request"})
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(key, value)
		}
	}
	c.Writer.WriteHeader(resp.StatusCode)
}
