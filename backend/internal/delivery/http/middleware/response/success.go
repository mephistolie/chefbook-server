package response

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body"
	"net/http"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func NewId(c *gin.Context, id int, message string) {
	Success(c, response_body.Id{Id: id, Message: message})
}

func Message(c *gin.Context, message string) {
	Success(c, response_body.Message{Message: message})
}

func Link(c *gin.Context, link string) {
	Success(c, response_body.Link{Link: link})
}
