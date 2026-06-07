package service

import (
	"auth-api/model"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user model.User) (model.User, error)
}

type AuthService struct {
	userRepository UserRepository
}

func NewAuthService(userRepository UserRepository) *AuthService {
	return &AuthService{
		userRepository: userRepository,
	}
}

func (s *AuthService) Signup(email string, password string) (model.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	return s.userRepository.Create(user)
}
