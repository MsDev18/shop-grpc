package router

func (r Router) registerHealthRoute () {
	r.engine.GET("/health-check" , r.healthHandler.HealthCheck)
}