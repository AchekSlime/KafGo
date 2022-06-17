package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

func newWarnResponse(c *gin.Context, statusCode int, message string) {
	logrus.Warning(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}
