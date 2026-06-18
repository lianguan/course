package v1

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ultrathreads/internal/domain"
	"ultrathreads/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"

	studentCtx = "studentId"
	adminCtx   = "adminId"
	userCtx    = "userId"
	schoolCtx  = "school"
	domainCtx  = "domain"
)

func (h *Handler) setSchoolFromRequest(c *gin.Context) {
	host := parseRequestHost(c)

	school, err := h.services.Schools.GetByDomain(c.Request.Context(), host)
	if err != nil {
		logger.Error(err)
		c.AbortWithStatus(http.StatusForbidden)

		return
	}

	c.Set(schoolCtx, school)
	c.Set(domainCtx, host)
}

func parseRequestHost(c *gin.Context) string {
	refererHeader := c.Request.Header.Get("Referer")
	refererParts := strings.Split(refererHeader, "/")

	// this logic is used to avoid crashes during integration testing
	if len(refererParts) < 3 {
		return c.Request.Host
	}

	hostParts := strings.Split(refererParts[2], ":")

	return hostParts[0]
}

func getSchoolFromContext(c *gin.Context) (domain.School, error) {
	value, ex := c.Get(schoolCtx)
	if !ex {
		return domain.School{}, errors.New("school is missing from ctx")
	}

	school, ok := value.(domain.School)
	if !ok {
		return domain.School{}, errors.New("failed to convert value from ctx to domain.School")
	}

	return school, nil
}

func (h *Handler) studentIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(studentCtx, id)
}

func (h *Handler) adminIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(adminCtx, id)
}

func (h *Handler) userIdentity(c *gin.Context) {
	id, err := h.parseAuthHeader(c)
	if err != nil {
		newResponse(c, http.StatusUnauthorized, err.Error())
		return
	}

	c.Set(userCtx, id)
}

func (h *Handler) parseAuthHeader(c *gin.Context) (string, error) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) == 1 {
		// 支持不带 Bearer 前缀的格式
		if len(headerParts[0]) == 0 {
			return "", errors.New("token is empty")
		}
		return h.tokenManager.Parse(headerParts[0])
	}
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func getStudentId(c *gin.Context) (uint, error) {
	return getIdByContext(c, studentCtx)
}

func getUserId(c *gin.Context) (uint, error) {
	return getIdByContext(c, userCtx)
}

func getIdByContext(c *gin.Context, context string) (uint, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return 0, errors.New("studentCtx not found")
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return 0, errors.New("studentCtx is of invalid type")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func getDomainFromContext(c *gin.Context) (string, error) {
	val, ex := c.Get(domainCtx)
	if !ex {
		return "", errors.New("domainCtx not found")
	}

	valStr, ok := val.(string)
	if !ok {
		return "", errors.New("domainCtx is of invalid type")
	}

	return valStr, nil
}
