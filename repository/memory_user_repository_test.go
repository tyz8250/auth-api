package repository

import (
	"auth-api/model"
	"testing"
)

func TestMemoryUserRepositoryCreateSavesUser(t *testing.T) {
	repo := NewMemoryUserRepository()

	saved, err := repo.Create(model.User{
		Email:        "user@example.com",
		PasswordHash: "hashed_password",
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}

	if saved.ID == 0 {
		t.Fatal("expected saved user to have ID")
	}

	found, ok := repo.FindByID(saved.ID)
	if !ok {
		t.Fatal("expected saved user to be found")
	}

	if found.Email != "user@example.com" {
		t.Fatalf("expected email user@example.com, got %q", found.Email)
	}

	if found.PasswordHash != "hashed_password" {
		t.Fatalf("expected password hash to be saved, got %q", found.PasswordHash)
	}
}
