package router

import "shop/internal/entity"

func (r Router) registerCategoryRoute() {
	categoryG := r.engine.Group("/category")

	categoryG.POST("", r.authMiddleware.AuthRequired(), r.authMiddleware.RoleRequired(entity.AdminRole), r.categoryHandler.Create)
	categoryG.PATCH("/:slug", r.authMiddleware.AuthRequired(), r.authMiddleware.RoleRequired(entity.AdminRole), r.categoryHandler.Update)
	categoryG.DELETE("/:slug", r.authMiddleware.AuthRequired(), r.authMiddleware.RoleRequired(entity.AdminRole), r.categoryHandler.Delete)
	
	categoryG.GET("", r.categoryHandler.GetAll)
	categoryG.GET("/:slug", r.categoryHandler.GetOne)
}
