package images

import (
	"context"
)

func NewDummyGenerator(img []byte) Generator {
	return dummyGenerator{img: img}
}

type dummyGenerator struct {
	img []byte
}

func (d dummyGenerator) Next(ctx context.Context) []byte {
	return d.img
}
