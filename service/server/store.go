package server

import (
	"io"
	"uni-token-service/store"

	"github.com/gin-gonic/gin"
	"go.etcd.io/bbolt"
)

func SetupStoreAPI(router gin.IRouter) {
	api := router.Group("/store").Use(RequireUserLogin())
	{
		api.DELETE("/:name", handleStoreDeleteAll)
		api.GET("/:name/:key", handleStoreGet)
		api.POST("/:name/:key", handleStorePut)
		api.DELETE("/:name/:key", handleStoreDelete)
	}
}

var createdBuckets = map[string]bool{}

func ensureBucket(name string) error {
	if created, exists := createdBuckets[name]; exists && created {
		return nil
	}
	return store.Db.Update(func(tx *bbolt.Tx) error {
		createdBuckets[name] = true
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		return err
	})
}

func handleStoreDeleteAll(c *gin.Context) {
	name := c.Param("name")
	err := store.DeleteBucket(name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

func handleStoreGet(c *gin.Context) {
	name := c.Param("name")
	key := c.Param("key")
	err := ensureBucket(name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = store.Db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(name))
		v := b.Get([]byte(key))
		if v == nil {
			c.JSON(200, nil)
			return nil
		}
		c.Data(200, "application/json", v)
		return nil
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
}

func handleStorePut(c *gin.Context) {
	name := c.Param("name")
	key := c.Param("key")
	err := ensureBucket(name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = store.Db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(name))
		return b.Put([]byte(key), body)
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}

func handleStoreDelete(c *gin.Context) {
	name := c.Param("name")
	key := c.Param("key")
	err := ensureBucket(name)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	err = store.Db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(name))
		return b.Delete([]byte(key))
	})
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.Status(200)
}
