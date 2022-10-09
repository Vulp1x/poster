package images

import (
	"bytes"
	"context"
	"fmt"
	"image"
)

func NewNoiseGenerator(imgBytes []byte, gamma float32) (Generator, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	return noiseGenerator{img: img, gammaCoef: gamma}, nil
}

type noiseGenerator struct {
	img       image.Image
	gammaCoef float32
}

func (g noiseGenerator) Next(ctx context.Context) []byte {
	// imageRect := g.img.Bounds().Size()
	// // result := noise.Generate(imageRect.X, imageRect.Y, &noise.Options{Monochrome: false, NoiseFn: noise.Uniform})
	//
	// buf := &bytes.Buffer{}
	//
	// err := jpeg.Encode(buf, dst, nil)
	// if err != nil {
	// 	logger.Errorf(ctx, "failed to encode new image: %v", err)
	// 	return nil
	// }
	//
	// return buf.Bytes()
	return nil
}
