package delayer

import (
	"math"
	"testing"
	"time"
)

func generateDelays(base time.Duration, max time.Duration, factor float64) []time.Duration {
	delay := base
	delays := []time.Duration{}
	for delay >= base && delay < max {
		delays = append(delays, delay)
		delay *= time.Duration(factor)
	}
	return delays
}

func TestExponentialBackoff(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		base   time.Duration
		max    time.Duration
		factor float64
	}{
		{time.Second, 128 * time.Second, 2},
		{3 * time.Second, 256 * time.Second, 2},
		{time.Second, 16 * time.Second, 0.5},
	}

	for _, tc := range testCases {
		delays := generateDelays(tc.base, tc.max, tc.factor)
		for i := 0; i < len(delays); i++ {
			delay := computeExponentialBackoffDelay(tc.base, tc.max, tc.factor, int16(i))
			if delay != delays[i] {
				t.Errorf("expected %d-th delay to be %s, got: %s", i, delays[i], delay)
			}
		}

		overflowingDelay := computeExponentialBackoffDelay(tc.base, tc.max, tc.factor, int16(len(delays)))
		if tc.factor > 1 {
			if overflowingDelay != tc.max {
				t.Errorf("expected delay to bounded by and equal to %s, got: %s", tc.max, overflowingDelay)
			}
		} else {
			if overflowingDelay != tc.base {
				t.Errorf("expected delay to be equal to %s, got: %s", tc.base, overflowingDelay)
			}
		}
	}
}

func TestOverflow(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		base    time.Duration
		max     time.Duration
		factor  float64
		attempt int16
	}{
		{time.Second, 32 * time.Second, 2, 5},
		{time.Second, 225 * time.Second, 15, 2},
		{time.Second, 128 * time.Second, 2, 7},
		{time.Second, 128 * time.Second, 2, 8},
		{time.Second, 128 * time.Second, 2, 1024},
		{time.Second, 128 * time.Second, 3, 5},
		{time.Second, 128 * time.Second, 3, math.MaxInt16},
		{time.Second, 128 * time.Second, math.MaxFloat64, math.MaxInt16},
	}

	for _, tc := range testCases {
		delay := computeExponentialBackoffDelay(tc.base, tc.max, tc.factor, tc.attempt)
		if delay != tc.max {
			t.Errorf("overflow error, expected max=[%v], got [%v]", tc.max, delay)
		}
	}
}
