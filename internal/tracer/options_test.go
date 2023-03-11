package tracer_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/inst-api/poster/internal/dbmodel"
	"github.com/inst-api/poster/internal/tracer"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func BenchmarkWithCustomAttributes(b *testing.B) {
	status := dbmodel.DoneTaskStatus
	id := uuid.New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tracer.WithCustomStruct("params", dbmodel.UpdateTaskStatusParams{
			Status: status,
			ID:     id,
		})
	}
}

func TestWithCustomAttributes(t *testing.T) {
	type args struct {
		name  string
		model interface{}
	}
	tests := []struct {
		name string
		args args
		want trace.SpanStartEventOption
	}{

		{
			name: "OK",
			args: args{
				name:  "oki",
				model: dbmodel.UpdateTaskParams{},
			},
			want: trace.WithAttributes(attribute.String("oki", "{TextTemplate: Title: Images:[] AccountNames:[] AccountLastNames:[] AccountUrls:[] AccountProfileImages:[] LandingAccounts:[] FollowTargets:false NeedPhotoTags:false PerPostSleepSeconds:0 PhotoTagsDelaySeconds:0 PostsPerBot:0 TargetsPerPost:0 ID:00000000-0000-0000-0000-000000000000 PhotoTargetsPerPost:0 PhotoTagsPostsPerBot:0 FixedTag:<nil> FixedPhotoTag:<nil>}")),
		},
	}

	opts := cmpopts.IgnoreUnexported(attribute.KeyValue{}, attribute.Value{})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println()
			if got := tracer.WithCustomStruct(tt.args.name, tt.args.model); !cmp.Equal(got, tt.want, opts) {
				t.Errorf("WithCustomAttributes() = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want, opts))
			}
		})
	}
}
