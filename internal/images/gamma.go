package images

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"math/rand"

	"github.com/disintegration/gift"
	"github.com/inst-api/poster/pkg/logger"
)

func NewGammaGenerator(imgBytes []byte, gamma float32) (Generator, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return gammaGenerator{img: img, gammaCoef: gamma}, nil
}

type gammaGenerator struct {
	img       image.Image
	gammaCoef float32
}

func (g gammaGenerator) Next(ctx context.Context) []byte {
	filter := gift.New(gift.Gamma(g.gammaCoef))
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

func NewRandomGammaGenerator(imgBytes []byte) (Generator, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return randomGammaGenerator{img: img}, nil
}

type randomGammaGenerator struct {
	img image.Image
}

func (g randomGammaGenerator) Next(ctx context.Context) []byte {
	filter := gift.New(gift.Gamma(0.4 + rand.Float32()))
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
