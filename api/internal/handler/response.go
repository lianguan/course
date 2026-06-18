package handler

import (
	"github.com/gin-gonic/gin"
	"ultrathreads/pkg/logger"
)

// DataResponse 数据响应结构
type DataResponse struct {
	Data  interface{} `json:"data"`
	Count int64       `json:"count"`
}

// IdResponse ID 响应结构
type IdResponse struct {
	ID interface{} `json:"id"`
}

// Response 通用响应结构
type Response struct {
	Message string `json:"message"`
}

// NewResponse 创建响应并记录日志
func NewResponse(c *gin.Context, statusCode int, message string) {
	// 区分客户端错误(4xx)和服务端错误(5xx)的日志级别
	if statusCode >= 500 {
		logger.Error(message)
	} else {
		logger.Warn(message)
	}
	c.AbortWithStatusJSON(statusCode, Response{message})
}
