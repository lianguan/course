package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ultrathreads/internal/domain"
	"ultrathreads/pkg/logger"

	"github.com/gin-gonic/gin"
)

const (
	AuthorizationHeader = "Authorization"
	SchoolCtx           = "school"
	DomainCtx           = "domain"
)

// SchoolsService 学校服务接口（避免循环依赖）
type SchoolsService interface {
	GetByDomain(ctx context.Context, domainName string) (domain.School, error)
}

// ParseRequestHost 从请求中解析域名
func ParseRequestHost(c *gin.Context) string {
	refererHeader := c.Request.Header.Get("Referer")
	refererParts := strings.Split(refererHeader, "/")

	// this logic is used to avoid crashes during integration testing
	if len(refererParts) < 3 {
		return c.Request.Host
	}

	hostParts := strings.Split(refererParts[2], ":")

	return hostParts[0]
}

// GetSchoolFromContext 从上下文获取学校信息
func GetSchoolFromContext(c *gin.Context) (domain.School, error) {
	value, ex := c.Get(SchoolCtx)
	if !ex {
		return domain.School{}, errors.New("school is missing from ctx")
	}

	school, ok := value.(domain.School)
	if !ok {
		return domain.School{}, errors.New("failed to convert value from ctx to domain.School")
	}

	return school, nil
}

// GetIdByContext 从上下文获取 ID
func GetIdByContext(c *gin.Context, context string) (uint, error) {
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

// GetDomainFromContext 从上下文获取域名
func GetDomainFromContext(c *gin.Context) (string, error) {
	val, ex := c.Get(DomainCtx)
	if !ex {
		return "", errors.New("domainCtx not found")
	}

	valStr, ok := val.(string)
	if !ok {
		return "", errors.New("domainCtx is of invalid type")
	}

	return valStr, nil
}

// TokenParser Token 解析接口（避免循环依赖）
type TokenParser interface {
	Parse(token string) (string, error)
}

// ParseAuthHeader 解析认证头
func ParseAuthHeader(c *gin.Context, tokenManager TokenParser) (string, error) {
	header := c.GetHeader(AuthorizationHeader)
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) == 1 {
		if len(headerParts[0]) == 0 {
			return "", errors.New("token is empty")
		}
		return tokenManager.Parse(headerParts[0])
	}
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return tokenManager.Parse(headerParts[1])
}

// ParseIdFromPath 从路径解析 ID
func ParseIdFromPath(c *gin.Context, param string) (uint, error) {
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

// SetSchoolFromRequest 从请求设置学校信息到上下文
func SetSchoolFromRequest(schoolsService SchoolsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		host := ParseRequestHost(c)

		school, err := schoolsService.GetByDomain(c.Request.Context(), host)
		if err != nil {
			logger.Error(err)
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Set(SchoolCtx, school)
		c.Set(DomainCtx, host)
	}
}
