package imageprocessor

const (
	UPLOADS_ROOT = "./uploads"
	UPLOADS_URL_PREFIX = "/uploads"
)

type Config struct {
	MaxDimension int `koanf:"max_dimension"`
}