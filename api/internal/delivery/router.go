package delivery

import (
	"fmt"
	"net/http"

	"ultrathreads/docs"
	"ultrathreads/internal/config"
	"ultrathreads/internal/delivery/admin"
	"ultrathreads/internal/delivery/public"
	"ultrathreads/internal/delivery/student"
	"ultrathreads/internal/service"
	"ultrathreads/pkg/auth"
	"ultrathreads/pkg/limiter"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		services:     services,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
		corsMiddleware,
	)

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if cfg.Environment != config.EnvLocal {
		docs.SwaggerInfo.Host = cfg.HTTP.Host
	}

	if cfg.Environment != config.EnvProd {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.InitRoutes(router.Group("/api"))

	return router
}

// InitRoutes registers all API routes under the given router group.
func (h *Handler) InitRoutes(api *gin.RouterGroup) {
	v1 := api.Group("/v1")

	adminHandler := admin.NewHandler(h.services, h.tokenManager)
	adminHandler.Init(v1)

	studentHandler := student.NewHandler(h.services, h.tokenManager)
	studentHandler.Init(v1)

	publicHandler := public.NewHandler(h.services, h.tokenManager)
	publicHandler.Init(v1)
}
