package router

import (
	_ "chefbook-server/docs"
	"chefbook-server/internal/app/dependencies/service"
	"chefbook-server/internal/config"
	"chefbook-server/internal/delivery/http/middleware"
	"chefbook-server/internal/delivery/http/router/v1"
	"chefbook-server/pkg/auth"
	"chefbook-server/pkg/limiter"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Router struct {
	services       *service.Service
	authMiddleware middleware.AuthMiddleware
	fileMiddleware middleware.FileMiddleware
}

func NewRouter(services *service.Service, tokenManager auth.TokenManager) *Router {
	authMiddleware := middleware.NewAuth(tokenManager)
	fileMiddleware := middleware.NewFile()
	return &Router{
		services:       services,
		authMiddleware: *authMiddleware,
		fileMiddleware: *fileMiddleware,
	}
}

func (r *Router) Init(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.Environment)

	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
	)

	r.initAPI(router)

	return router
}

func (r *Router) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewV1Router(r.services, r.authMiddleware, r.fileMiddleware)
	api := router.Group("/")
	{
		if gin.Mode() != gin.ReleaseMode {
			router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}

		handlerV1.Init(api)
	}
}
