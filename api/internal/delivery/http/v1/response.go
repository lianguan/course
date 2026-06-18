package v1

import (
	"github.com/gin-gonic/gin"
	"ultrathreads/pkg/logger"
)

type dataResponse struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}

type idResponse struct {
	ID interface{} `json:"id"`
}

type response struct {
	Message string `json:"message"`
}

func newResponse(c *gin.Context, statusCode int, message string) {
	// 区分客户端错误(4xx)和服务端错误(5xx)的日志级别
	if statusCode >= 500 {
		logger.Error(message)
	} else {
		logger.Warn(message)
	}
	c.AbortWithStatusJSON(statusCode, response{message})
}
