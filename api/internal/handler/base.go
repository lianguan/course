package handler

import (
	"ultrathreads/internal/service"
	"ultrathreads/pkg/auth"
)

// BaseHandler 基础处理器，包含公共依赖和方法
type BaseHandler struct {
	services     *service.Services
	tokenManager auth.TokenManager
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler(services *service.Services, tokenManager auth.TokenManager) *BaseHandler {
	return &BaseHandler{
		services:     services,
		tokenManager: tokenManager,
	}
}

// Services 获取服务层
func (h *BaseHandler) Services() *service.Services {
	return h.services
}

// TokenManager 获取 Token 管理器
func (h *BaseHandler) TokenManager() auth.TokenManager {
	return h.tokenManager
}
