package delayer

import (
	"time"
)

// Delayer рассчитывает по номеру попытки задержку перед следующей попыткой выполнить задачу.
type Delayer func(attempt int16) time.Duration

// NewExpBackoffDelayer - реализация Jitter Exponential Backoff.
// Подробно про алгоритм: https://aws.amazon.com/ru/blogs/architecture/exponential-backoff-and-jitter.
//
// jitter - [FullJitter | EqualJitter | ZeroJitter];
// задает нижнюю границу диапазона, в котором рандомизируется задержка,
// верхней границей является значение задержки,
//
// base - начальное значение задержки,
//
// max - максимальное значение задержки,
//
// factor - число, на степень которого умножается base.
func NewExpBackoffDelayer(jitter jitter, base, max time.Duration, factor float64) Delayer {
	return func(attempt int16) time.Duration {
		delay := computeExponentialBackoffDelay(base, max, factor, attempt)
		return time.Duration(jitter.applyTo(float64(delay)))
	}
}

// NewJitterDelayer рассчитывает случайную задержку в диапазоне, зависящим от jitter.
// Используйте FullJitter, EqualJitter или ZeroJitter.
func NewJitterDelayer(jitter jitter, max time.Duration) Delayer {
	return func(_ int16) time.Duration {
		return time.Duration(jitter.applyTo(float64(max)))
	}
}
