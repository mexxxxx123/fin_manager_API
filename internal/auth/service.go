package auth

import (
	"errors"
	"fin_manager_API/m/internal/user"
	"fin_manager_API/m/pkg/di"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository di.IUserRepository
}

func NewAuthService(userRepo di.IUserRepository) *AuthService {
	return &AuthService{UserRepository: userRepo}
}

func (service *AuthService) Login(email string, password string) (string, error) {
	user, _ := service.UserRepository.FindByEmail(email)
	if user == nil {
		return "", errors.New(ErrDontFindAcc)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrDontFindPass)
	}
	return user.Email, nil
}

func (service *AuthService) Register(email, password, name string) (string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return "", errors.New(ErrUserExist)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return "", err
	}
	return user.Email, err
}
