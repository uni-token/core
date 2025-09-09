package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"uni-token-service/constants"
	"uni-token-service/logic"
)

func SetupActionAPI(router *gin.Engine) {
	router.GET("/", handleCheck)
	router.POST("/ui/active", handleUIActive)
}

func handleCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"__uni_token": true,
		"version":     constants.Version,
	})
}

func handleUIActive(c *gin.Context) {
	var req struct {
		Session string `json:"session"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	shouldContinue := logic.OnUIActive(req.Session)

	c.JSON(http.StatusOK, gin.H{"continue": shouldContinue})
}
