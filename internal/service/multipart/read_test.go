package multipart

import (
	"context"
	"os"
	"testing"
)

func TestName(t *testing.T) {
	f, err := os.Open("testdata/100 акк.txt")
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	defer f.Close()

	bots, err := readBotsList(context.Background(), f)
	if err != nil {
		t.Fatalf("got errors: %v", err)
	}

	if len(bots) != 100 {
		t.Fatalf("got %d bots instead of 100", len(bots))
	}

}
