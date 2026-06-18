package public

import (
	"ultrathreads/internal/handler"

	"github.com/gin-gonic/gin"
)

func (h *Handler) setSchoolFromRequest(c *gin.Context) {
	handler.SetSchoolFromRequest(h.Services().Schools)(c)
}

func getSchoolFromContext(c *gin.Context) (interface{}, error) {
	return handler.GetSchoolFromContext(c)
}

func getDomainFromContext(c *gin.Context) (string, error) {
	return handler.GetDomainFromContext(c)
}
