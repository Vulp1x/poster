package images

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/h2non/bimg"
)

func TestNewNoiseGenerator(t *testing.T) {
	type args struct {
		imgBytes []byte
		gamma    float32
	}
	tests := []struct {
		name    string
		args    args
		want    Generator
		wantErr bool
	}{
		{
			name: "Ok",
			args: args{
				imgBytes: testImageBytes,
				gamma:    0.1,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewNoiseGenerator(tt.args.imgBytes, tt.args.gamma)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGammaGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			f, err := os.Create(fmt.Sprintf("test_cat_noise_%f.jpeg", tt.args.gamma))
			if err != nil {
				log.Fatalf("os.Create failed: %v", err)
			}

			nextImageBytes := got.Next(context.Background())

			f.Write(nextImageBytes)
			f.Close()

		})
	}
}

func TestMetadata(t *testing.T) {
	meta, err := bimg.Metadata(testImageBytes)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(meta)

}
