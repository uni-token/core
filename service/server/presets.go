package server

import (
	"net/http"
	"sort"
	"time"
	"uni-token-service/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddPresetRequest struct {
	Name      string   `json:"name" binding:"required"`
	Providers []string `json:"providers"`
}

type UpdatePresetRequest struct {
	Name      string   `json:"name" binding:"required"`
	Providers []string `json:"providers"`
}

// Helper function to check if preset name exists (excluding specified ID)
func isPresetNameExists(name string, excludeID string) (bool, error) {
	presets, err := store.AppPresets.List()
	if err != nil {
		return false, err
	}

	for _, preset := range presets {
		if preset.Name == name && preset.ID != excludeID {
			return true, nil
		}
	}
	return false, nil
}

func SetupPresetsAPI(router gin.IRouter) {
	api := router.Group("/presets").Use(RequireUserLogin())
	{
		api.GET("/list", handlePresetsList)
		api.POST("/add", handleAddPreset)
		api.DELETE("/delete/:id", handleDeletePreset)
		api.PUT("/update/:id", handleUpdatePreset)
	}
}

func handlePresetsList(c *gin.Context) {
	presets, err := store.AppPresets.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve presets"})
		return
	}

	// Sort presets by UpdatedAt in descending order (most recent first)
	sort.Slice(presets, func(i, j int) bool {
		return presets[i].UpdatedAt.After(presets[j].UpdatedAt)
	})

	c.JSON(http.StatusOK, gin.H{"data": presets})
}

func handleAddPreset(c *gin.Context) {
	var req AddPresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if preset name already exists
	exists, err := isPresetNameExists(req.Name, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check preset name"})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Preset name already exists"})
		return
	}

	now := time.Now()
	presetID := uuid.New().String()
	preset := store.AppPreset{
		ID:        presetID,
		Name:      req.Name,
		Providers: req.Providers,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = store.AppPresets.Put(presetID, preset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save preset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preset added successfully", "data": preset})
}

func handleDeletePreset(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Preset ID is required"})
		return
	}

	// Get preset to check if it's the default preset
	preset, err := store.AppPresets.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Preset not found"})
		return
	}

	if preset.Name == "default" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete default preset"})
		return
	}

	err = store.AppPresets.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete preset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preset deleted successfully"})
}

func handleUpdatePreset(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Preset ID is required"})
		return
	}

	var req UpdatePresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get existing preset to preserve CreatedAt and current values
	existingPreset, err := store.AppPresets.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Preset not found"})
		return
	}

	// Check if new name conflicts with other presets (excluding current one)
	if req.Name != existingPreset.Name {
		exists, err := isPresetNameExists(req.Name, id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check preset name"})
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Preset name already exists"})
			return
		}
	}

	preset := store.AppPreset{
		ID:        existingPreset.ID,
		Name:      req.Name,
		Providers: req.Providers,
		CreatedAt: existingPreset.CreatedAt,
		UpdatedAt: time.Now(),
	}

	err = store.AppPresets.Put(id, preset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preset"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preset updated successfully", "data": preset})
}
