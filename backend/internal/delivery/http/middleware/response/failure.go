package response

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body"
	"net/http"
)

func Failure(c *gin.Context, err error)  {
	response := response_body.NewError(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, response)
}