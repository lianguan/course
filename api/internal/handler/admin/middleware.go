package admin

import (
	"net/http"

	"ultrathreads/internal/handler"

	"github.com/gin-gonic/gin"
)

const (
	adminCtx = "adminId"
)

func (h *Handler) setSchoolFromRequest(c *gin.Context) {
	handler.SetSchoolFromRequest(h.Services().Schools)(c)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	id, err := handler.ParseAuthHeader(c, h.TokenManager())
	if err != nil {
		handler.NewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(adminCtx, id)
}

func getAdminId(c *gin.Context) (uint, error) {
	return handler.GetIdByContext(c, adminCtx)
}

func getSchoolFromContext(c *gin.Context) (interface{}, error) {
	return handler.GetSchoolFromContext(c)
}

func getDomainFromContext(c *gin.Context) (string, error) {
	return handler.GetDomainFromContext(c)
}
