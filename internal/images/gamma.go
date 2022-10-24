package images

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"math/rand"

	"github.com/disintegration/gift"
	"github.com/inst-api/poster/internal/domain"
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

// NewRandomGammaGenerator создает генератор на основе набора картинок.
// Next применяет гамма коррекцию со случайной величиной к случайной картинки из набора
func NewRandomGammaGenerator(imagesBytes [][]byte) (Generator, error) {
	images := make([]image.Image, len(imagesBytes))
	for i, imageBytes := range imagesBytes {
		img, _, err := image.Decode(bytes.NewReader(imageBytes))
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %v", err)
		}

		images[i] = img
	}

	return randomGammaGenerator{images: images}, nil
}

type randomGammaGenerator struct {
	images []image.Image
}

func (g randomGammaGenerator) Next(ctx context.Context) []byte {
	randImage := domain.RandomFromSlice(g.images)

	filter := gift.New(gift.Gamma(0.8 + rand.Float32()*0.6))
	dst := image.NewRGBA(filter.Bounds(randImage.Bounds()))
	filter.Draw(dst, randImage)

	buf := &bytes.Buffer{}

	err := jpeg.Encode(buf, dst, &jpeg.Options{Quality: 60 + rand.Intn(40)})
	if err != nil {
		logger.Errorf(ctx, "failed to encode new image: %v", err)
		return nil
	}

	return buf.Bytes()
}
