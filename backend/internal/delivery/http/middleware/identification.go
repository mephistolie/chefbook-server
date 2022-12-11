package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/delivery/http/middleware/response"
	"github.com/mephistolie/chefbook-server/internal/entity/failure"
	"github.com/mephistolie/chefbook-server/pkg/auth"
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

func (m AuthMiddleware) GetUserId(c *gin.Context) (string, error) {
	return m.getIdByContext(c, userContext)
}

func (m AuthMiddleware) getIdByContext(c *gin.Context, context string) (string, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return "", failure.Unknown
	}

	id, ok := idFromCtx.(string)
	if !ok {
		return "", failure.Unknown
	}

	return id, nil
}
