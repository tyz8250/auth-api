package service

import (
	"auth-api/model"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type spyUserRepository struct {
	savedUser model.User
}

func (r *spyUserRepository) Create(user model.User) (model.User, error) {
	r.savedUser = user

	user.ID = 1
	user.CreatedAt = "2025-10-15T12:00:00Z"
	user.UpdatedAt = "2025-10-15T12:00:00Z"

	return user, nil
}

func TestAuthServiceSignupHashesPassword(t *testing.T) {
	userRepository := &spyUserRepository{}
	authService := NewAuthService(userRepository)

	user, err := authService.Signup("user@example.com", "secret")
	if err != nil {
		t.Fatalf("Signup returned error: %v", err)
	}

	if user.Email != "user@example.com" {
		t.Fatalf("expected email user@example.com, got %q", user.Email)
	}

	if user.PasswordHash == "" {
		t.Fatal("expected password hash to be set")
	}

	if user.PasswordHash == "secret" {
		t.Fatal("password hash should not be the plain password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte("secret")); err != nil {
		t.Fatalf("password hash should match original password: %v", err)
	}

	if userRepository.savedUser.PasswordHash == "" {
		t.Fatal("expected password hash to be saved through repository")
	}

	if userRepository.savedUser.PasswordHash == "secret" {
		t.Fatal("repository should not receive the plain password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userRepository.savedUser.PasswordHash), []byte("secret")); err != nil {
		t.Fatalf("repository should receive a hash that matches original password: %v", err)
	}
}
