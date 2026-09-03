package router

import "shop/internal/entity"

func (r Router) registerProductRoute() {
	productG := r.engine.Group("/product")

	productG.POST("", r.authMiddleware.AuthRequired(), r.authMiddleware.RoleRequired(entity.AdminRole), r.productHandler.Create)
	productG.GET("/:slug" , r.productHandler.GetOneBySlug)
	productG.GET("" , r.productHandler.GetAll)
}
