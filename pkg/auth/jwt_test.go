package auth

import (
	"testing"
)

const testSecret = "test-secret-key"

func TestNewTokenAndParseToken(t *testing.T) {
	token, err := NewToken(42, testSecret, 24)
	if err != nil {
		t.Fatalf("NewToken: %v", err)
	}
	if token == "" {
		t.Error("expected non-empty token")
	}
	claims, err := ParseToken(token, testSecret)
	if err != nil {
		t.Fatalf("ParseToken: %v", err)
	}
	if claims.AuthorID != 42 {
		t.Errorf("AuthorID want 42, got %d", claims.AuthorID)
	}
}

func TestParseToken_InvalidSecret(t *testing.T) {
	token, _ := NewToken(1, testSecret, 1)
	_, err := ParseToken(token, "wrong-secret")
	if err == nil {
		t.Error("ParseToken with wrong secret should fail")
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	_, err := ParseToken("not.a.token", testSecret)
	if err == nil {
		t.Error("ParseToken with invalid token should fail")
	}
}
