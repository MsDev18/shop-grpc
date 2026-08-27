package router

import (
	"shop/internal/api/handler/address"
	"shop/internal/api/handler/auth"
	"shop/internal/api/handler/category"
	"shop/internal/api/handler/health"
	"shop/internal/api/handler/product"
	"shop/internal/api/handler/province"
	"shop/internal/api/handler/user"
	authmiddleware "shop/internal/api/middleware/auth"

	"github.com/gin-gonic/gin"
)

type Router struct {
	engine *gin.Engine
	// handlers statements
	healthHandler   health.Handler
	authHandler     auth.Handler
	userHandler     user.Handler
	categoryHandler category.Handler
	provinceHandler province.Handler
	addressHandler  address.Handler
	productHandler  product.Handler
	authMiddleware  authmiddleware.Middleware
}

func New(engine *gin.Engine, healthHandler health.Handler, authHandler auth.Handler, userHandler user.Handler, categoryHandler category.Handler, provinceHandler province.Handler, addressHandler address.Handler, productHandler product.Handler, authMiddleware authmiddleware.Middleware) Router {
	return Router{
		engine: engine,
		// handlers statements
		healthHandler:   healthHandler,
		authHandler:     authHandler,
		userHandler:     userHandler,
		categoryHandler: categoryHandler,
		provinceHandler: provinceHandler,
		addressHandler:  addressHandler,
		productHandler:  productHandler,
		authMiddleware:  authMiddleware,
	}
}

func (r Router) Register() {
	r.registerHealthRoute()
	r.registerAuthRoute()
	r.registerUserRoute()
	r.registerStaticRoute()
	r.registerCategoryRoute()
	r.registerProvinceRoute()
	r.registerAddressRoute()
	r.registerProductRoute()
}
