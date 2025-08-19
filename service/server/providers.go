package server

import (
	"net/http"
	"uni-token-service/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddProviderRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
	BaseURL  string `json:"baseUrl" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

func SetupProvidersAPI(router gin.IRouter) {
	api := router.Group("/providers").Use(RequireUserLogin())
	{
		api.GET("/list", handleProvidersList)
		api.POST("/add", handleAddProvider)
		api.DELETE("/delete/:id", handleDeleteProvider)
		api.PUT("/update/:id", handleUpdateProvider)
	}
}

func handleProvidersList(c *gin.Context) {
	providers, err := store.LLMProviders.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve providers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": providers})
}

func handleAddProvider(c *gin.Context) {
	var req AddProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID based on type
	providerID := uuid.New().String()

	provider := store.LLMProvider{
		ID:       providerID,
		Name:     req.Name,
		Type:     req.Type,
		Protocol: req.Protocol,
		BaseURL:  req.BaseURL,
		Token:    req.Token,
	}

	err := store.LLMProviders.Put(providerID, provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider added successfully", "data": provider})
}

func handleDeleteProvider(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider ID is required"})
		return
	}

	// Check if provider exists
	_, err := store.LLMProviders.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	// First, remove this provider from all presets
	presets, err := store.AppPresets.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve presets"})
		return
	}

	for _, preset := range presets {
		var updatedProviders []string
		modified := false

		for _, p := range preset.Providers {
			if p != id {
				updatedProviders = append(updatedProviders, p)
			} else {
				modified = true
			}
		}

		if modified {
			preset.Providers = updatedProviders
			err := store.AppPresets.Put(preset.ID, preset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update presets"})
				return
			}
		}
	}

	// Then delete the provider
	err = store.LLMProviders.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider deleted successfully"})
}

func handleUpdateProvider(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider ID is required"})
		return
	}

	var req AddProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing provider to check if it exists
	existingProvider, err := store.LLMProviders.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	provider := store.LLMProvider{
		ID:       existingProvider.ID, // Keep the original ID
		Name:     req.Name,
		Type:     req.Type,
		Protocol: req.Protocol,
		BaseURL:  req.BaseURL,
		Token:    req.Token,
	}

	err = store.LLMProviders.Put(id, provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider updated successfully", "data": provider})
}
