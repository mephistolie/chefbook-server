package v1

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/models"
	"net/http"
	"strconv"
)

const MaxUploadSize = 1 << 20
var ImageTypes = map[string]interface{} {
	"image/jpeg": nil,
	"image/png": nil,
}

func (h *Handler) initUsersRoutes(api *gin.RouterGroup) {
	auth := api.Group("/users", h.userIdentity)
	{
		auth.GET("", h.getUserInfo)
		auth.POST("/change-name", h.setUserName)
		auth.POST("/avatar", h.uploadAvatar)
		auth.DELETE("/avatar", h.deleteAvatar)
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
			Broccoins: user.Broccoins,
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

func (h *Handler) setUserName(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var username models.UserNameInput
	if err := c.BindJSON(&username); err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidInput.Error())
		return
	}

	err = h.services.SetUserName(userId, username.Username)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespUsernameChanged,
	})
}

func (h *Handler) uploadAvatar(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		newResponse(c, http.StatusBadRequest, models.ErrInvalidFileInput.Error())
		return
	}

	defer func() {
		err := file.Close()
		if err != nil {
			newResponse(c, http.StatusInternalServerError, err.Error())
			return
		}
	}()

	buffer := make([]byte, fileHeader.Size)
	_, err = file.Read(buffer)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	fileType := http.DetectContentType(buffer)

	fileBytes := bytes.NewReader(buffer)

	if _, ex := ImageTypes[fileType]; !ex {
		newResponse(c, http.StatusBadRequest, models.ErrFileTypeNotSupported.Error())
		return
	}

	url, err := h.services.UploadAvatar(c.Request.Context(), userId, fileBytes, fileHeader.Size, "image/jpeg")
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"link": url,
	})
}

func (h *Handler) deleteAvatar(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.DeleteAvatar(c.Request.Context(), userId)
	if err != nil {
		newResponse(c, http.StatusInternalServerError, models.ErrUnableDeleteAvatar.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"message": models.RespAvatarDeleted,
	})
}