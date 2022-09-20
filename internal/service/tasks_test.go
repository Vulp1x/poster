package service

import (
	"fmt"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestName(t *testing.T) {
	hashPass, err := bcrypt.GenerateFromPassword([]byte("admin0"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	fmt.Println(string(hashPass))

	err = bcrypt.CompareHashAndPassword(hashPass, []byte("admin0"))
	fmt.Println(err)
}
