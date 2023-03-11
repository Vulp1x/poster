package tracer

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// WithCustomStruct добавляет модель к трейсу
func WithCustomStruct(name string, model interface{}) trace.SpanStartEventOption {
	return trace.WithAttributes(attribute.Stringer(name, stringer{model: model}))
}

type stringer struct {
	model interface{}
}

func (s stringer) String() string {
	return fmt.Sprintf("%+v\n", s.model)
}
