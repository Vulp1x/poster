package images

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/disintegration/gift"
	"github.com/inst-api/poster/pkg/logger"
)

func NewNoiseGenerator(imgBytes []byte, gamma float32) (Generator, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return convolutionalGenerator{img: img, gammaCoef: gamma}, nil
}

type convolutionalGenerator struct {
	img       image.Image
	gammaCoef float32
}

func (g convolutionalGenerator) Next(ctx context.Context) []byte {
	filter := gift.New(gift.Convolution([]float32{-1, 1, 2, 0, 0, 2, 1, 1, 2}, true, false, false, 0))

	dst := image.NewRGBA(filter.Bounds(g.img.Bounds()))
	filter.Draw(dst, g.img)

	buf := &bytes.Buffer{}

	err := jpeg.Encode(buf, dst, nil)
	if err != nil {
		logger.Errorf(ctx, "failed to encode new image: %v", err)
		return nil
	}

	return buf.Bytes()
}
