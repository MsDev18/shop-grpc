package server

type Env string

const (
	EnvDevelopment Env = "development"
	EnvProduction  Env = "production"
)

func (e Env) IsValid() bool {
	switch e {
	case EnvDevelopment, EnvProduction:
		return true
	default:
		return false
	}
}
