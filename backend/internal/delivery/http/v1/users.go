package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	auth := api.Group("/users")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
		auth.POST("/sign-out", h.signOut)
		auth.GET("/activate/:link", h.activate)
		auth.POST("/refresh", h.refreshSession)
		auth.GET("/", h.userIdentity, h.getUserInfo)
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

	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := h.services.SignOut(userId, input.RefreshToken); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
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

func (h *Handler) getUserInfo(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	userProfileId, err := strconv.Atoi(c.Request.URL.Query().Get("user_id"))
	if err != nil {
		userProfileId = userId
	}

	user, err := h.services.GetUserInfo(userProfileId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userId == user.Id {
		c.JSON(http.StatusOK, models.UserDetailedInfo{
			Id:        user.Id,
			Username:  user.Username.String,
			Email:     user.Email,
			Avatar:    user.Avatar.String,
			Premium:   user.Premium.Time,
			IsBlocked: user.IsBlocked,
		})
	} else {
		c.JSON(http.StatusOK, models.UserInfo{
			Id:       user.Id,
			Username: user.Username.String,
			Avatar:   user.Avatar.String,
			Premium:  user.Premium.Time,
		})
	}
}
