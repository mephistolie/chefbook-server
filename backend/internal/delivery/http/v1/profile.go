package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/model"
	"net/http"
	"strconv"
)

func (h *Handler) getUserInfo(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	userProfileId, err := strconv.Atoi(c.Request.URL.Query().Get("user_id"))
	if err != nil {
		userProfileId = userId
	}

	user, err := h.services.GetUserInfo(userProfileId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userId == user.Id {
		c.JSON(http.StatusOK, model.UserDetailedInfo{
			Id:                user.Id,
			Username:          user.Username.String,
			Email:             user.Email,
			CreationTimestamp: user.CreationTimestamp,
			Avatar:            user.Avatar.String,
			Premium:           user.Premium.Time,
			Broccoins:         user.Broccoins,
			IsBlocked:         user.IsBlocked,
		})
	} else {
		c.JSON(http.StatusOK, model.UserInfo{
			Id:       user.Id,
			Username: user.Username.String,
			Avatar:   user.Avatar.String,
			Premium:  user.Premium.Time,
		})
	}
}

func (h *Handler) setUserName(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var username model.UserNameInput
	if err := c.BindJSON(&username); err != nil {
		newErrorResponse(c, http.StatusBadRequest, model.ErrInvalidInput.Error())
		return
	}

	err = h.services.Profile.SetUsername(userId, username.Username)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newMessageResponse(c, RespUsernameChanged)
}

func (h *Handler) getUserKey(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	url, err := h.services.Encryption.GetUserKeyLink(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newLinkResponse(c, url)
}

func (h *Handler) uploadAvatar(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	file, err := getFileByCtx(c)
	if err != nil {
		return
	}

	if _, ex := ImageTypes[file.ContentType]; !ex {
		newErrorResponse(c, http.StatusBadRequest, model.ErrFileTypeNotSupported.Error())
		return
	}

	url, err := h.services.Profile.UploadAvatar(c.Request.Context(), userId, file)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	newLinkResponse(c, url)
}

func (h *Handler) deleteAvatar(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Profile.DeleteAvatar(c.Request.Context(), userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, model.ErrUnableDeleteAvatar.Error())
		return
	}

	newMessageResponse(c, RespAvatarDeleted)
}