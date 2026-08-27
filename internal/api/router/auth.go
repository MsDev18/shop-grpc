package router

func (r Router) registerAuthRoute() {
	authG := r.engine.Group("/auth")

	authG.POST("/send-otp", r.authHandler.SendOtp)
	authG.POST("/check-otp", r.authHandler.CheckOtp)
	authG.GET("/me", r.authMiddleware.AuthRequired(), r.authHandler.Me)
	authG.POST("/logout", r.authMiddleware.AuthRequired(), r.authHandler.Logout)
	authG.POST("/refresh-token", r.authHandler.RefreshToken)
}
