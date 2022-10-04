package requests

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/inst-api/poster/internal/domain"
)

func TestPrepareUploadRequest(t *testing.T) {
	type args struct {
		b     domain.BotAccount
		image []byte
	}
	tests := []struct {
		name string
		args args
		want *http.Request
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrepareUploadRequest(tt.args.b, tt.args.image); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrepareUploadRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
