package server

import (
	"net/http"
	"strconv"
	"uni-token-service/store"

	"github.com/gin-gonic/gin"
)

// SetupUsageAPI sets up usage-related API endpoints
func SetupUsageAPI(router gin.IRouter) {
	api := router.Group("/usage").Use(RequireUserLogin())
	{
		api.GET("/stats", handleGetUsageStats)
		api.GET("/list", handleGetUsageList)
		api.POST("/clear", handleClearUsageRecords)
	}
}

// handleGetUsageStats returns aggregated usage statistics
func handleGetUsageStats(c *gin.Context) {
	daysParam := c.DefaultQuery("days", "30")
	days, err := strconv.Atoi(daysParam)
	if err != nil || days <= 0 {
		days = 30
	}

	stats, err := store.GetUsageStats(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get usage stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}

// handleGetUsageList returns paginated usage records
func handleGetUsageList(c *gin.Context) {
	usages, err := store.Usage.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get usage list"})
		return
	}

	// Simple pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 200 {
		pageSize = 50
	}

	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(usages) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data": gin.H{
				"records":  []interface{}{},
				"total":    len(usages),
				"page":     page,
				"pageSize": pageSize,
			},
		})
		return
	}

	if end > len(usages) {
		end = len(usages)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"records":  usages[start:end],
			"total":    len(usages),
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// handleClearUsageRecords clears all usage records
func handleClearUsageRecords(c *gin.Context) {
	err := store.Usage.Clear()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear usage records"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Usage records cleared successfully",
	})
}
