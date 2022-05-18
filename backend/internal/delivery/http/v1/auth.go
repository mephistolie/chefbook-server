package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/model"
	"net/http"
)

func (h *Handler) signUp(c *gin.Context) {
	var input model.AuthData

	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, model.ErrInvalidInput)
		return
	}

	id, err := h.services.Auth.SignUp(input)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	newIdResponse(c, id, RespActivationLink)
}

func (h *Handler) activate(c *gin.Context) {
	activationLink, err := uuid.Parse(c.Param("link"))
	if err != nil {
		newErrorResponse(c, model.ErrInvalidInput)
		return
	}

	if err := h.services.Auth.ActivateUser(activationLink); err != nil {
		newErrorResponse(c, model.ErrInvalidInput)
		return
	}

	newMessageResponse(c, RespProfileActivated)
}

func (h *Handler) signIn(c *gin.Context) {
	var input model.AuthData
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, model.ErrInvalidInput)
		return
	}

	res, err := h.services.Auth.SignIn(input, c.Request.RemoteAddr)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) signOut(c *gin.Context) {
	var input model.RefreshInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, model.ErrInvalidInput)
		return
	}

	if err := h.services.Auth.SignOut(input.RefreshToken); err != nil {
		newErrorResponse(c, err)
		return
	}

	newMessageResponse(c, RespSignOutSuccessfully)
}

func (h *Handler) refreshSession(c *gin.Context) {
	var input model.RefreshInput
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, model.ErrInvalidInput)
		return
	}

	res, err := h.services.Auth.RefreshSession(input.RefreshToken, c.Request.RemoteAddr)
	if err != nil {
		newErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, res)
}
