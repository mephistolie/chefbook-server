package handler

import (
	"chefbook-server/internal/app/dependencies/service"
	"chefbook-server/internal/delivery/http/middleware/response"
	"chefbook-server/internal/delivery/http/presentation/request_body"
	"chefbook-server/internal/delivery/http/presentation/response_body"
	"chefbook-server/internal/delivery/http/presentation/response_body/message"
	"chefbook-server/internal/entity/failure"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ParamActivationCode = "activation_code"
)

type AuthHandler struct {
	service service.Auth
}

func NewAuthHandler(service service.Auth) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// SignUp Swagger Documentation
// @Summary Sign Up
// @Tags auth
// @Description Create new profile
// @Accept json
// @Produce json
// @Param input body request_body.Credentials true "Credentials"
// @Success 200 {object} response_body.Id
// @Failure 400 {object} response_body.Error
// @Router /v1/auth/sign-up [post]
func (h *AuthHandler) SignUp(c *gin.Context) {
	var body request_body.Credentials
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}
	if err := body.Validate(); err != nil {
		response.Failure(c, err)
		return
	}

	id, err := h.service.SignUp(body.Entity())
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.NewId(c, id, message.ActivationLinkSent)
}

// ActivateProfile Swagger Documentation
// @Summary Activate Profile
// @Tags auth
// @Description Activate profile
// @Accept json
// @Produce json
// @Param activation_code path string true "Activation Code"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/auth/activate/{activation_code} [get]
func (h *AuthHandler) ActivateProfile(c *gin.Context) {
	activationLink, err := uuid.Parse(c.Param(ParamActivationCode))
	if err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	if err := h.service.ActivateProfile(activationLink); err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.ProfileActivated)
}

// SignIn Swagger Documentation
// @Summary Sign In
// @Tags auth
// @Description Sign in to profile
// @Accept json
// @Produce json
// @Param input body request_body.Credentials true "Credentials"
// @Success 200 {object} response_body.Tokens
// @Failure 400 {object} response_body.Error
// @Router /v1/auth/sign-in [post]
func (h *AuthHandler) SignIn(c *gin.Context) {
	var body request_body.Credentials
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	tokens, err := h.service.SignIn(body.Entity(), c.Request.RemoteAddr)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, response_body.NewTokens(tokens))
}

// SignOut Swagger Documentation
// @Summary Sign Out
// @Tags auth
// @Description Sign out and remove session
// @Accept json
// @Produce json
// @Param input body request_body.RefreshToken true "Refresh Token"
// @Success 200 {object} response_body.Message
// @Failure 400 {object} response_body.Error
// @Router /v1/auth/sign-out [post]
func (h *AuthHandler) SignOut(c *gin.Context) {
	var body request_body.RefreshToken
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	if err := h.service.SignOut(body.RefreshToken); err != nil {
		response.Failure(c, err)
		return
	}

	response.Message(c, message.SignOutSuccessfully)
}

// RefreshSession Swagger Documentation
// @Summary Refresh Session
// @Tags auth
// @Description Refresh session to get new tokens pair
// @Accept json
// @Produce json
// @Param input body request_body.RefreshToken true "Refresh Token"
// @Success 200 {object} response_body.Tokens
// @Failure 400 {object} response_body.Error
// @Router /v1/auth/refresh [post]
func (h *AuthHandler) RefreshSession(c *gin.Context) {
	var body request_body.RefreshToken
	if err := c.BindJSON(&body); err != nil {
		response.Failure(c, failure.InvalidBody)
		return
	}

	tokens, err := h.service.RefreshSession(body.RefreshToken, c.Request.RemoteAddr)
	if err != nil {
		response.Failure(c, err)
		return
	}

	response.Success(c, response_body.NewTokens(tokens))
}
