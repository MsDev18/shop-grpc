package config

import (
	"shop/internal/api/server"
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
	Server      server.Config         `koanf:"server"`
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
