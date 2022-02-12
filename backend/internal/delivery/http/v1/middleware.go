package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/model"
	"net/http"
	"strconv"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx = "userId"
)

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", model.ErrEmptyAuthHeader
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", model.ErrInvalidAuthHeader
	}

	if len(headerParts[1]) == 0 {
		return "", model.ErrEmptyToken
	}

	return h.tokenManager.Parse(headerParts[1])
}

func (h *Handler) userIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	c.Set(userCtx, id)
}

func getUserId(c *gin.Context) (int, error) {
	return getIdByContext(c, userCtx)
}

func getIdByContext(c *gin.Context, context string) (int, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return -1, model.ErrUserIdNotFound
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return -1, model.ErrInvalidUserId
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1, err
	}

	return id, nil
}