package public

import (
	"ultrathreads/internal/handler"
	"ultrathreads/internal/service"
	"ultrathreads/pkg/auth"

	"github.com/gin-gonic/gin"
)

// Handler 公开处理器
type Handler struct {
	*handler.BaseHandler
}

// NewHandler 创建公开处理器
func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		BaseHandler: handler.NewBaseHandler(services, tokenManager),
	}
}

// Init 初始化公开路由
func (h *Handler) Init(api *gin.RouterGroup) {
	h.initCoursesRoutes(api)
	h.initCallbackRoutes(api)

	api.GET("/settings", h.setSchoolFromRequest, h.getSchoolSettings)
	api.GET("/promocodes/:code", h.setSchoolFromRequest, h.getPromo)
	api.GET("/offers/:id", h.setSchoolFromRequest, h.getOffer)
}
