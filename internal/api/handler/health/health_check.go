package health

import (
	"shop/internal/pkg/response"

	"github.com/gin-gonic/gin"
)


func (h Handler) HealthCheck (ctx *gin.Context) {
	response.New(ctx).OK("Health Check ✅", nil)
}