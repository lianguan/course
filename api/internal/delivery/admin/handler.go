package admin

import (
	"errors"
	"strconv"

	"ultrathreads/internal/service"
	"ultrathreads/pkg/auth"

	"github.com/gin-gonic/gin"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	h.initAdminRoutes(api)
}

func parseIdFromPath(c *gin.Context, param string) (uint, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return 0, errors.New("empty id param")
	}

	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, errors.New("invalid id param")
	}

	if id == 0 {
		return 0, errors.New("invalid id param")
	}

	return uint(id), nil
}
