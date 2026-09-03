package router

func (r Router) registerUserRoute() {
	userG := r.engine.Group("/user")

	userG.GET("/profile", r.authMiddleware.AuthRequired(), r.userHandler.Profile)
	userG.PATCH("/update-profile", r.authMiddleware.AuthRequired(), r.userHandler.UpdateProfile)
	userG.PATCH("/change-password", r.authMiddleware.AuthRequired(), r.userHandler.ChangePassword)
}
