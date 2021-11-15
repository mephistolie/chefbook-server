package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/service"
	"github.com/mephistolie/chefbook-server/pkg/auth"
)

type Handler struct {
	services     *service.Service
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Service, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services: services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup)  {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initRecipesRoutes(v1)
	}
}