package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-out", h.signOut)
		auth.GET("/activate/:link", h.activate)
		auth.POST("/refresh", h.refreshSession)
	}
}

func (h *Handler) signUp(c *gin.Context) {
	var input models.AuthData

	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	id, err := h.services.Users.SignUp(input)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id":      id,
		"message": models.RespActivationLink,
	})
}

func (h *Handler) activate(c *gin.Context) {
	activationLink, err := uuid.Parse(c.Param("link"))
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	if err := h.services.Users.ActivateUser(activationLink); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespProfileActivated,
	})
}

func (h *Handler) signIn(c *gin.Context) {
	var input models.AuthData
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	res, err := h.services.SignIn(input, c.Request.RemoteAddr)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) || errors.Is(err, models.ErrAuthentication) {
			newResponse(c, http.StatusBadRequest, models.ErrAuthentication.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) signOut(c *gin.Context) {
	var input models.RefreshInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	if err := h.services.SignOut(input.RefreshToken); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespSignOutSuccessfully,
	})
}

func (h *Handler) refreshSession(c *gin.Context) {
	var input models.RefreshInput
	if err := c.BindJSON(&input); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	res, err := h.services.RefreshSession(input.RefreshToken, c.Request.RemoteAddr)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) || errors.Is(err, models.ErrAuthentication) {
			newResponse(c, http.StatusBadRequest, models.ErrAuthentication.Error())
			return
		}
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, res)
}