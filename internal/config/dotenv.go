package config

import (
	"fmt"
	"strings"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf/providers/confmap"
)

func (k KoanfConfig) LoadFromDotEnv(path string) {
	// Read .env File 
	dotEnvVar, err := godotenv.Read(path)
	if err != nil {
		panic(fmt.Errorf("Error loading .env file: %v", err))
	}
	// Map the env keys 
	var envVar = make(map[string]any)
	for k , v := range dotEnvVar {
		newKey := strings.ReplaceAll(strings.ToLower(k), "__" , ".")
		envVar[newKey] = v
	}
	// load envVar in koanf core 
	k.koanf.Load(confmap.Provider(envVar , "."),nil)
}
