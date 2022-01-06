package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/mephistolie/chefbook-server/internal/service"
	"github.com/mephistolie/chefbook-server/pkg/auth"
)

const MaxUploadSize = 1 << 20
var ImageTypes = map[string]interface{} {
	"image/jpeg": nil,
	"image/png": nil,
}

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
		h.initAuthRoutes(v1)
		h.initUsersRoutes(v1)
		h.initRecipesRoutes(v1)
		h.initCategoriesRoutes(v1)
		h.initShoppingListRoutes(v1)
	}
}