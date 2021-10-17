package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/siruspen/logrus"
)

type Error struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, message string)  {
	logrus.Errorf(message)
	c.AbortWithStatusJSON(statusCode, Error{message})
}
