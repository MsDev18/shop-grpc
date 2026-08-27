package main

import "shop/internal/config"


func LoadConfig() config.Config {
	appConfig := config.New()
	appConfig.LoadFromDotEnv(".env")
	appConfig.LoadFromYml("config.yml")
	return appConfig.GetConfig()
}