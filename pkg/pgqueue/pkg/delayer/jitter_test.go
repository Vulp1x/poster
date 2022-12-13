package delayer

import (
	"testing"
)

func TestJitter(t *testing.T) {
	t.Parallel()
	jitters := []jitter{FullJitter, EqualJitter, ZeroJitter}
	values := []float64{1, 2, 4, 8, 3, 9, 27, 81}
	for _, j := range jitters {
		for _, value := range values {
			jittered := j.applyTo(value)
			lowerBound := value * (1.0 - j.value)
			if jittered < lowerBound || jittered > value {
				t.Errorf("expected result in the range [%f, %f], got: %f", lowerBound, value, jittered)
			}
		}
	}
}
