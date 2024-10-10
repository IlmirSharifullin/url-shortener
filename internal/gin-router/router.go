package gin_router

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
	"log/slog"
	"net/http"
	"url-shortener/internal/storage"
)

func InitEngine(logger *slog.Logger, storage storage.Storage) (r *gin.Engine) {
	r = gin.New()
	r.Use(sloggin.New(logger))

	r.Use(gin.Recovery())
	r.Use(DBMiddleware(&storage))

	r.POST("/add", addAlias)

	r.GET("/*any", func(c *gin.Context) {
		path := c.Param("any")
		if path == "/ping" {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		} else {
			getByAlias(c)
		}
	})

	return
}

func DBMiddleware(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("storage", *db)
		c.Next()
	}
}

func addAlias(c *gin.Context) {
	var db storage.Storage
	_db, ok := c.Get("storage")
	db = _db.(storage.Storage)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal",
		})
	}
	var json struct {
		Url   string `json:"url"`
		Alias string `json:"alias"`
	}

	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong json object",
			"error":   err,
		})
	}

	err := db.SaveURL(json.Url, json.Alias)
	if err != nil {
		if errors.Is(err, storage.ErrUrlExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": fmt.Sprintf("%v", err),
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Some error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("%s added to storage with %s alias", json.Url, json.Alias),
	})
}

func getByAlias(c *gin.Context) {
	var db storage.Storage
	_db, ok := c.Get("storage")

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
			"url":     "",
		})
		return
	}

	db, ok = _db.(storage.Storage) // Assert to the interface, not the struct
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal error",
			"url":     "",
		})
		return
	}

	alias := c.Request.URL.Path[1:]
	url, err := db.GetURL(alias)

	if err != nil {
		if errors.Is(err, storage.ErrNoRow) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "no match for this alias",
				"url":     "",
			})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "unknown error",
				"url":     "",
			})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}
