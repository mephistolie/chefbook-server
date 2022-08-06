package middleware

import (
	"chefbook-server/internal/delivery/http/middleware/response"
	"chefbook-server/internal/entity/failure"
	"chefbook-server/pkg/auth"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userContext         = "userId"
)

type AuthMiddleware struct {
	tokenManager auth.TokenManager
}

func NewAuth(manager auth.TokenManager) *AuthMiddleware {
	return &AuthMiddleware{
		tokenManager: manager,
	}
}

func (m AuthMiddleware) CheckUserIdentity(c *gin.Context) {
	id, err := m.parseAuthHeader(c)
	if err != nil {
		response.Failure(c, err)
	}
	c.Set(userContext, id)
}

func (m AuthMiddleware) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", failure.EmptyAuthHeader
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", failure.InvalidAuthHeader
	}

	if len(headerParts[1]) == 0 {
		return "", failure.EmptyToken
	}

	token, err := m.tokenManager.Parse(headerParts[1])
	if err != nil {
		return "", failure.InvalidToken
	}
	return token, err
}

func (m AuthMiddleware) GetUserId(c *gin.Context) (int, error) {
	return m.getIdByContext(c, userContext)
}

func (m AuthMiddleware) getIdByContext(c *gin.Context, context string) (int, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return 0, failure.Unknown
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return 0, failure.Unknown
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, failure.Unknown
	}

	return id, nil
}
