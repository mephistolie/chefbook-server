package response

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body"
)

func Failure(c *gin.Context, err error) {
	statusCode, response := response_body.NewError(err)
	c.AbortWithStatusJSON(statusCode, response)
}
