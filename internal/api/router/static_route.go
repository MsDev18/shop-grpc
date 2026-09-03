package router

import (
	"github.com/gin-gonic/gin"
)

func (r Router) registerStaticRoute() {
	r.engine.StaticFS("/uploads", gin.Dir("./uploads", false))
}
