package server

import (
	"net/http"
	"net/url"
	"time"

	"uni-token-service/logic"
	"uni-token-service/store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupAppAPI(router gin.IRouter) {
	router.POST("/app/register", handleAppRegister)
	router.Group("/app").Use(RequireUserLogin()).GET("/list", handleAppList)
	router.Group("/app").Use(RequireUserLogin()).POST("/toggle", handleAppToggle)
	router.Group("/app").Use(RequireUserLogin()).DELETE("/delete/:id", handleAppDelete)
	router.Group("/app").Use(RequireUserLogin()).DELETE("/clear", handleAppClear)
}

var waitForGrant = make(map[string]chan bool)

func handleAppRegister(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description" binding:"required"`
		UID         string `json:"uid"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := req.UID

	granted := func() {
		c.JSON(http.StatusOK, gin.H{
			"token": uid,
		})
	}

	info := store.AppInfo{
		ID:           uid,
		Name:         req.Name,
		Description:  req.Description,
		Provider:     "",
		Granted:      false,
		CreatedAt:    time.Now(),
		LastActiveAt: time.Now(),
	}

	if uid != "" {
		app, err := store.Apps.Get(uid)
		if err == nil {
			if app.Granted {
				updateLastActiveTime(uid)
				granted()
				return
			}
			if app.Provider != "" {
				info.Provider = app.Provider
			}
			if app.CreatedAt != (time.Time{}) {
				info.CreatedAt = app.CreatedAt
			}
		}
	} else {
		// Search for same-name app
		apps, err := store.Apps.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve apps"})
			return
		}
		for _, app := range apps {
			if app.Name == req.Name {
				uid = app.ID
				if app.Provider != "" {
					info.Provider = app.Provider
				}
				if app.CreatedAt != (time.Time{}) {
					info.CreatedAt = app.CreatedAt
				}
				break
			}
		}

		// Else, generate a new UID
		if uid == "" {
			uid = uuid.NewString()
		}
		info.ID = uid
	}

	store.Apps.Put(uid, info)

	channel := make(chan bool)
	waitForGrant[uid] = channel

	params := url.Values{
		"action":         {"register"},
		"appId":          {uid},
		"appName":        {req.Name},
		"appDescription": {req.Description},
	}
	logic.OpenUI(params)

	select {
	case result := <-channel:
		if result {
			granted()
		} else {
			c.JSON(http.StatusForbidden, gin.H{"error": "App registration denied"})
		}
	case <-time.After(60 * time.Second):
		c.JSON(http.StatusRequestTimeout, gin.H{"error": "App registration timed out"})
	}
	delete(waitForGrant, uid)
}

func handleAppList(c *gin.Context) {
	apps, err := store.Apps.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve apps"})
		return
	}

	c.JSON(http.StatusOK, apps)
}

func handleAppToggle(c *gin.Context) {
	var req struct {
		ID       string `json:"id" binding:"required"`
		Granted  bool   `json:"granted"`
		Provider string `json:"provider"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if info, err := store.Apps.Get(req.ID); err == nil {
		info.Granted = req.Granted
		if req.Provider != "" {
			info.Provider = req.Provider
		}
		store.Apps.Put(req.ID, info)
		action := "authorized"
		if !req.Granted {
			action = "revoked"
		}
		c.JSON(http.StatusOK, gin.H{"message": "App " + action + " successfully"})

		if waitForGrant[req.ID] != nil {
			waitForGrant[req.ID] <- req.Granted
		}
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
	}
}

func updateLastActiveTime(appID string) {
	if app, err := store.Apps.Get(appID); err == nil {
		app.LastActiveAt = time.Now()
		store.Apps.Put(appID, app)
	}
}

func handleAppDelete(c *gin.Context) {
	appID := c.Param("id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "App ID is required"})
		return
	}

	// Check if app exists
	_, err := store.Apps.Get(appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "App not found"})
		return
	}

	// Delete the app
	err = store.Apps.Delete(appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete app"})
		return
	}

	// Clean up any pending grant channel
	if waitForGrant[appID] != nil {
		close(waitForGrant[appID])
		delete(waitForGrant, appID)
	}

	c.JSON(http.StatusOK, gin.H{"message": "App deleted successfully"})
}

func handleAppClear(c *gin.Context) {
	// Get all apps first
	apps, err := store.Apps.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve apps"})
		return
	}

	// Delete all apps
	for _, app := range apps {
		err = store.Apps.Delete(app.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete some apps"})
			return
		}

		// Clean up any pending grant channel
		if waitForGrant[app.ID] != nil {
			close(waitForGrant[app.ID])
			delete(waitForGrant, app.ID)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "All apps deleted successfully"})
}
