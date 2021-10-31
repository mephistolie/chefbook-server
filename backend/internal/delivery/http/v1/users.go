package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	models2 "github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	auth := api.Group("/users")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-out", h.signIn)
		auth.GET("/activate/:link", h.activate)
		auth.GET("/refresh", h.refreshSession)
	}
}

func (h *Handler) signUp(c *gin.Context) {
	var input models2.AuthData

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, models2.ErrInvalidInput.Error())
		return
	}

	id, err := h.services.Users.SignUp(input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": models2.RespActivationLink,
	})
}

func (h *Handler) activate(c *gin.Context) {
	activationLink, err := uuid.Parse(c.Param("link"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models2.ErrInvalidInput.Error())
		return
	}

	if err := h.services.Users.ActivateUser(activationLink); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models2.RespProfileActivated,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input models2.AuthData
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, models2.ErrInvalidInput.Error())
		return
	}

	res, err := h.services.SignIn(input, c.Request.RemoteAddr)
	if err != nil {
		if errors.Is(err, models2.ErrUserNotFound) || errors.Is(err, models2.ErrAuthentication) {
			newResponse(c, http.StatusBadRequest, models2.ErrAuthentication.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) refreshSession(c *gin.Context) {
	var input models2.RefreshInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, models2.ErrInvalidInput.Error())
		return
	}

	res, err := h.services.RefreshSession(input.Token, c.Request.RemoteAddr)
	if err != nil {
		if errors.Is(err, models2.ErrUserNotFound) || errors.Is(err, models2.ErrAuthentication) {
			newResponse(c, http.StatusBadRequest, models2.ErrAuthentication.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}