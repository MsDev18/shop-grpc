package config

import (
	"shop/internal/pkg/imageprocessor"
	"shop/internal/repository/mysql"
	authservice "shop/internal/service/auth"

	"github.com/knadh/koanf/v2"
)

type KoanfConfig struct {
	koanf *koanf.Koanf
}

type Config struct {
	MySQL       mysql.Config          `koanf:"mysql"`
	AuthService authservice.Config    `koanf:"auth_service"`
	Upload imageprocessor.Config `koanf:"upload"`
}

func New() KoanfConfig {
	k := koanf.New(".")
	return KoanfConfig{
		koanf: k,
	}
}

func (k KoanfConfig) GetConfig() Config {
	var cfg Config
	k.koanf.Unmarshal("", &cfg)
	return cfg
}

func (k KoanfConfig) GetMySQLConfig() mysql.Config {
	var cfg mysql.Config
	k.koanf.Unmarshal("mysql", &cfg)
	return cfg
}