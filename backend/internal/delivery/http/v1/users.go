package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
	"strconv"
)

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	auth := api.Group("/users", h.userIdentity)
	{
		auth.GET("", h.getUserInfo)
		auth.PUT("/change-name", h.getUserInfo)
	}
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
