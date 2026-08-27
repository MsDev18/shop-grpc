package imageprocessor

import (
	"context"
	"io"
)

type Processor struct {
	config Config
	storage Storage
}

type Storage interface {
	Save(ctx context.Context , filename string, data io.Reader) (string , error)
}

func New (config Config , storage Storage) Processor {
	return Processor{
		config: config,
		storage: storage,
	}
}