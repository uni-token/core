package server

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"uni-token-service/constants"
	"uni-token-service/logic"
)

func SetupActionAPI(router *gin.Engine) {
	router.GET("/", handleCheck)
	router.GET("/ui/open", handleOpenUI)
	router.POST("/ui/active", handleUIActive)
}

func handleCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"__uni_token": true,
		"version":     constants.Version,
	})
}

func handleOpenUI(c *gin.Context) {
	_, cleanup, err := logic.OpenUI(url.Values{}, true)
	if cleanup != nil {
		cleanup()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open UI"})
		return
	}
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
