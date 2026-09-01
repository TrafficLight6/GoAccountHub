package middleware

import (
	"github.com/TrafficLight6/GoAccountHub/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ConfigInsertMiddleware(config config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", config)
		c.Next()
	}
}

func DBInsertMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
