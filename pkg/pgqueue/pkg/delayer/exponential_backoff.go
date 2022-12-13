package delayer

import (
	"math"
	"time"
)

func computeExponentialBackoffDelay(base time.Duration, max time.Duration, factor float64, attempt int16) time.Duration {
	// Overflow guard, checking if:
	// base * factor**attempt >= max <=> attempt * log(factor) >= log(max / base)
	if float64(attempt)*math.Log(factor) >= math.Log(float64(max/base)) {
		return max
	}

	backoff := base * time.Duration(math.Pow(factor, float64(attempt)))
	if backoff < base {
		backoff = base
	}

	return backoff
}
