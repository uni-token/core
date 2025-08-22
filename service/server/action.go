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
	router.GET("/open", handleOpenUrl)
}

func handleCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"__uni_token": true,
		"version":     constants.GetVersion(),
	})
}

func handleOpenUrl(c *gin.Context) {
	err := logic.OpenUI(url.Values{}, true)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open UI"})
		return
	}
}
