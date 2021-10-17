package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input models.User

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
}

func (h *Handler) signIn(c *gin.Context) {

}