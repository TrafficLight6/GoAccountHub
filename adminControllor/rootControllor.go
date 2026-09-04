package adminControllor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Welcome to GoAccountHub",
	})
}
