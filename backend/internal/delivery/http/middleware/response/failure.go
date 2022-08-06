package response

import (
	"chefbook-server/internal/delivery/http/presentation/response_body"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Failure(c *gin.Context, err error) {
	response := response_body.NewError(err)
	c.AbortWithStatusJSON(http.StatusBadRequest, response)
}
