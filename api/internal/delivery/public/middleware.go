package public

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
	schoolCtx = "school"
	domainCtx = "domain"
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

func getIdByContext(c *gin.Context, context string) (uint, error) {
	idFromCtx, ok := c.Get(context)
	if !ok {
		return 0, errors.New("context not found")
	}

	idStr, ok := idFromCtx.(string)
	if !ok {
		return 0, errors.New("context is of invalid type")
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
