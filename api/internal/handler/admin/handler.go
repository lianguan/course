package admin

import (
	"ultrathreads/internal/handler"
	"ultrathreads/internal/service"
	"ultrathreads/pkg/auth"

	"github.com/gin-gonic/gin"
)

// Handler 管理员处理器
type Handler struct {
	*handler.BaseHandler
}

// NewHandler 创建管理员处理器
func NewHandler(services *service.Services, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		BaseHandler: handler.NewBaseHandler(services, tokenManager),
	}
}

// Init 初始化管理员路由
func (h *Handler) Init(api *gin.RouterGroup) {
	h.initAdminRoutes(api)
}
