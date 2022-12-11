package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mephistolie/chefbook-server/internal/app/dependencies/service"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware/response"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/request_body"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/presentation/response_body/message"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
)

const (
	ParamUserId = "user_id"

	queryUserId = "user_id"

	maxAvatarSize = 1 << 20
)

type ProfileHandler struct {
	authMiddleware middleware.AuthMiddleware
	fileMiddleware middleware.FileMiddleware
	service        service.Profile
}

func NewProfileHandler(authMiddleware middleware.AuthMiddleware, fileMiddleware middleware.FileMiddleware, service service.Profile) *ProfileHandler {
	return &ProfileHandler{
		authMiddleware: authMiddleware,
		fileMiddleware: fileMiddleware,
		service:        service,
	}
}

// GetProfileInfo Swagger Documentation
// @Summary Get Profile Info
// @Security ApiKeyAuth
// @Tags profile
// @Description Get user profile info
// @Accept json
// @Produce json
// @Param user_id query string false "User ID"
// @Success 200 {object} response_body.DetailedProfileInfo
// @Success 200 {object} response_body.MinimalProfileInfo
// @Failure 400 {object} response_body.Error
// @Router /v1/profile [get]
func (r *ProfileHandler) GetProfileInfo(c *gin.Context) {
	userId, err := r.authMiddleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	requestedUserId, err := uuid.Parse(c.Param(queryUserId))
	if err != nil {
		requestedUserId = userId
	}

	profile, err := r.service.GetProfile(requestedUserId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	if userId == profile.Id {
		response.Success(c, response_body.NewDetailedProfileInfo(profile))
	} else {
		response.Success(c, response_body.NewMinimalProfileInfoByProfile(profile))
	}
}

// ChangePassword Swagger Documentation
// @Summary Change Password
// @Security ApiKeyAuth
// @Tags profile
// @Description Change profile password
// @Accept json
// @Produce json
// @Param input body request_body.PasswordChanging true "Password Changing"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/profile/password [put]
func (r *ProfileHandler) ChangePassword(c *gin.Context) {
	userId, err := r.authMiddleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var body request_body.PasswordChanging
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	if err := body.Validate(); err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.ChangePassword(userId, body.OldPassword, body.NewPassword)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.PasswordChanged)
}

// SetUsername Swagger Documentation
// @Summary Change Username
// @Security ApiKeyAuth
// @Tags profile
// @Description Change profile username
// @Accept json
// @Produce json
// @Param input body request_body.Username true "Username"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/profile/username [put]
func (r *ProfileHandler) SetUsername(c *gin.Context) {
	userId, err := r.authMiddleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	var body request_body.Username
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	err = r.service.SetUsername(userId, body.Username)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.UsernameChanged)
}

// UploadAvatar Swagger Documentation
// @Summary Upload avatar
// @Security ApiKeyAuth
// @Tags profile
// @Description Upload profile avatar
// @Accept mpfd
// @Produce json
// @Param file formData file true "Avatar File"
// @Success 200 {object} response_body.Link
// @Failure 400 {object} response_body.Error
// @Router /v1/profile/avatar [put]
func (r *ProfileHandler) UploadAvatar(c *gin.Context) {
	userId, err := r.authMiddleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	file, err := r.fileMiddleware.GetFileWithMaxSize(c, maxAvatarSize)
	if err != nil {
		response.Failure(c, err)
		return
	}

	if !file.IsImage() {
		response.Failure(c, failure.UnsupportedFileType)
		return
	}

	url, err := r.service.UploadAvatar(c.Request.Context(), userId, file)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Link(c, url)
}

// DeleteAvatar Swagger Documentation
// @Summary Delete Avatar
// @Security ApiKeyAuth
// @Tags profile
// @Description Delete profile avatar
// @Accept json
// @Produce json
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/profile/avatar [delete]
func (r *ProfileHandler) DeleteAvatar(c *gin.Context) {
	userId, err := r.authMiddleware.GetUserId(c)
	if err != nil {
		response.Failure(c, err)
		return
	}

	err = r.service.DeleteAvatar(c.Request.Context(), userId)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.AvatarDeleted)
}
