package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	db "backend/database"
)

func DBHealthCheck(c *gin.Context) {
	if err := db.DB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"database": "connected",
	})
}
