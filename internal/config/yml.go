package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)


func (k KoanfConfig) LoadFromYml(path string) {
	err := k.koanf.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		fmt.Println("Error loading YAML file:", err)
	}
}