package router

func (r Router) registerProvinceRoute() {
	provinceG := r.engine.Group("/province")
	// set routes
	provinceG.GET("", r.provinceHandler.GetAll)
	provinceG.GET("/:id", r.provinceHandler.GetOne)
}
