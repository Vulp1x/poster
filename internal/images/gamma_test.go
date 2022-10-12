package images

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"
	"testing"
)

//go:embed testdata/br77.jpg
var testImageBytes []byte

func TestNewGammaGenerator(t *testing.T) {
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
		{
			name: "Ok",
			args: args{
				imgBytes: testImageBytes,
				gamma:    0.3,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Ok",
			args: args{
				imgBytes: testImageBytes,
				gamma:    0.6,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Ok",
			args: args{
				imgBytes: testImageBytes,
				gamma:    0.9,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Ok",
			args: args{
				imgBytes: testImageBytes,
				gamma:    1.5,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewGammaGenerator(tt.args.imgBytes, tt.args.gamma)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGammaGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			f, err := os.Create(fmt.Sprintf("test_cat_gamma_%f.jpeg", tt.args.gamma))
			if err != nil {
				log.Fatalf("os.Create failed: %v", err)
			}

			nextImageBytes := got.Next(context.Background())

			f.Write(nextImageBytes)
			f.Close()

		})
	}
}

func TestNewRandomGammaGenerator(t *testing.T) {
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
		{name: "Ok", args: args{imgBytes: testImageBytes}},
		{name: "Ok", args: args{imgBytes: testImageBytes}},
		{name: "Ok", args: args{imgBytes: testImageBytes}},
		{name: "Ok", args: args{imgBytes: testImageBytes}},
		{name: "Ok", args: args{imgBytes: testImageBytes}},
	}
	for i, tt := range tests {
		i := i
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRandomGammaGenerator(tt.args.imgBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGammaGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			f, err := os.Create(fmt.Sprintf("test_cat_random_gamma_%d.jpeg", i))
			if err != nil {
				log.Fatalf("os.Create failed: %v", err)
			}

			nextImageBytes := got.Next(context.Background())

			f.Write(nextImageBytes)
			f.Close()

		})
	}
}

func BenchmarkName(b *testing.B) {
	g, _ := NewGammaGenerator(testImageBytes, 0.7)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		g.Next(ctx)
	}
}
