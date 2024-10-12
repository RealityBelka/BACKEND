package handlers

import "github.com/gin-gonic/gin"

type BiometryHandlers interface {
	CreateFaceBiometry(c *gin.Context)
	CreateVoiceBiometry(c *gin.Context)
}
