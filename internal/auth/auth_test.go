package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestHashPassword(t *testing.T) {
	password := "testPassword123"
	hashed, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hashed == "" {
		t.Error("HashPassword returned empty string")
	}
	if hashed == password {
		t.Error("HashPassword returned plain text password")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	password := "testPassword123"
	hashed, _ := HashPassword(password)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{"correct password", password, hashed, true},
		{"wrong password", "wrongPassword", hashed, false},
		{"empty password", "", hashed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckPasswordHash(tt.password, tt.hash)
			if got != tt.want {
				t.Errorf("CheckPasswordHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret123"
	expiresIn := time.Hour

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}
	if token == "" {
		t.Error("MakeJWT returned empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "secret123"
	expiresIn := time.Hour

	token, _ := MakeJWT(userID, tokenSecret, expiresIn)

	tests := []struct {
		name    string
		token   string
		secret  string
		wantID  uuid.UUID
		wantErr bool
	}{
		{"valid token", token, tokenSecret, userID, false},
		{"invalid secret", token, "wrongSecret", uuid.Nil, true},
		{"empty token", "", tokenSecret, uuid.Nil, true},
		{"invalid token", "invalid.token.here", tokenSecret, uuid.Nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ValidateJWT(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.wantID {
				t.Errorf("ValidateJWT() = %v, want %v", got, tt.wantID)
			}
		})
	}
}

func TestMakeRefreshToken(t *testing.T) {
	token1, err := MakeRefreshToken()
	if err != nil {
		t.Fatalf("MakeRefreshToken failed: %v", err)
	}
	if token1 == "" {
		t.Error("MakeRefreshToken returned empty token")
	}

	token2, _ := MakeRefreshToken()
	if token1 == token2 {
		t.Error("MakeRefreshToken generated duplicate tokens")
	}
}
