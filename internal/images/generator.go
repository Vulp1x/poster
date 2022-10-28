package images

import (
	"context"
)

// Generator создает новую уникальную картинку из оригинала
type Generator interface {
	Next(ctx context.Context) []byte
}
