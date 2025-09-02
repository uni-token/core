package server

import (
	"net/http"
	"uni-token-service/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddKeyRequest struct {
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
	BaseURL  string `json:"baseUrl" binding:"required"`
	Token    string `json:"token" binding:"required"`
}

func SetupKeysAPI(router gin.IRouter) {
	api := router.Group("/keys").Use(RequireUserLogin())
	{
		api.GET("/list", handleKeysList)
		api.POST("/add", handleAddKey)
		api.DELETE("/delete/:id", handleDeleteKey)
		api.PUT("/update/:id", handleUpdateKey)
	}
}

func handleKeysList(c *gin.Context) {
	keys, err := store.LLMKeys.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve keys"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": keys})
}

func handleAddKey(c *gin.Context) {
	var req AddKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate ID based on type
	keyID := uuid.New().String()

	key := store.LLMKey{
		ID:       keyID,
		Name:     req.Name,
		Type:     req.Type,
		Protocol: req.Protocol,
		BaseURL:  req.BaseURL,
		Token:    req.Token,
	}

	err := store.LLMKeys.Put(keyID, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Key added successfully", "data": key})
}

func handleDeleteKey(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key ID is required"})
		return
	}

	// Check if key exists
	_, err := store.LLMKeys.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	// First, remove this key from all presets
	presets, err := store.AppPresets.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve presets"})
		return
	}

	for _, preset := range presets {
		var updatedKeys []string
		modified := false

		for _, p := range preset.Keys {
			if p != id {
				updatedKeys = append(updatedKeys, p)
			} else {
				modified = true
			}
		}

		if modified {
			preset.Keys = updatedKeys
			err := store.AppPresets.Put(preset.ID, preset)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update presets"})
				return
			}
		}
	}

	// Then delete the key
	err = store.LLMKeys.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Key deleted successfully"})
}

func handleUpdateKey(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Key ID is required"})
		return
	}

	var req AddKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing key to check if it exists
	existingKey, err := store.LLMKeys.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
		return
	}

	key := store.LLMKey{
		ID:       existingKey.ID, // Keep the original ID
		Name:     req.Name,
		Type:     req.Type,
		Protocol: req.Protocol,
		BaseURL:  req.BaseURL,
		Token:    req.Token,
	}

	err = store.LLMKeys.Put(id, key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update key"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Key updated successfully", "data": key})
}
