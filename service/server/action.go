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
	router.POST("/ui/opened", handleUIOpened)
}

func handleCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"__uni_token": true,
		"version":     constants.GetVersion(),
	})
}

func handleOpenUI(c *gin.Context) {
	err := logic.OpenUI(url.Values{}, true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open UI"})
		return
	}
}

func handleUIOpened(c *gin.Context) {
	var req struct {
		Session string `json:"session"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	logic.OnUIOpened(req.Session)
	c.Status(http.StatusOK)
}
