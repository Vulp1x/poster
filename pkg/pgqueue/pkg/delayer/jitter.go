package delayer

import (
	"math/rand"
)

type jitter struct {
	value float64
}

var (
	// Рандомизирует value в диапазоне [0, value].
	FullJitter = jitter{1.0}
	// Рандомизирует value в диапазоне [value/2, value].
	EqualJitter = jitter{0.5}
	// Сохраняет исходное значение.
	ZeroJitter = jitter{0.0}
)

// applyTo рандомизирует переданное значение.
func (j jitter) applyTo(value float64) float64 {
	return ((1.0 - j.value) + j.value*rand.Float64()) * value
}
