package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("secret123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if hash == "" || hash == "secret123" {
		t.Error("expected non-empty hash different from password")
	}
}

func TestCheckPassword(t *testing.T) {
	hash, _ := HashPassword("mypass")
	if !CheckPassword(hash, "mypass") {
		t.Error("CheckPassword should succeed for correct password")
	}
	if CheckPassword(hash, "wrong") {
		t.Error("CheckPassword should fail for wrong password")
	}
	if CheckPassword("", "any") {
		t.Error("CheckPassword should fail for empty hash")
	}
}
