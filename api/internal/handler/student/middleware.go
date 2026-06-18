package student

import (
	"net/http"

	"ultrathreads/internal/handler"

	"github.com/gin-gonic/gin"
)

const (
	studentCtx = "studentId"
	userCtx    = "userId"
)

func (h *Handler) setSchoolFromRequest(c *gin.Context) {
	handler.SetSchoolFromRequest(h.Services().Schools)(c)
}

func (h *Handler) studentIdentity(c *gin.Context) {
	id, err := handler.ParseAuthHeader(c, h.TokenManager())
	if err != nil {
		handler.NewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(studentCtx, id)
}

func (h *Handler) userIdentity(c *gin.Context) {
	id, err := handler.ParseAuthHeader(c, h.TokenManager())
	if err != nil {
		handler.NewResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, id)
}

func getStudentId(c *gin.Context) (uint, error) {
	return handler.GetIdByContext(c, studentCtx)
}

func getUserId(c *gin.Context) (uint, error) {
	return handler.GetIdByContext(c, userCtx)
}

func getSchoolFromContext(c *gin.Context) (interface{}, error) {
	return handler.GetSchoolFromContext(c)
}

func getDomainFromContext(c *gin.Context) (string, error) {
	return handler.GetDomainFromContext(c)
}
