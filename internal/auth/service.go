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

func (service *AuthService) Login(email string, password string) (uint, string, error) {
	user, _ := service.UserRepository.FindByEmail(email)
	if user == nil {
		return 0, "", errors.New(ErrDontFindPass)
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return 0, "", errors.New(ErrDontFindPass)
	}
	return user.ID, user.Email, nil
}

func (service *AuthService) Register(email, password, name string) (uint, string, error) {
	existedUser, _ := service.UserRepository.FindByEmail(email)
	if existedUser != nil {
		return 0, "", errors.New(ErrUserExist)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, "", err
	}

	user := &user.User{
		Email:    email,
		Password: string(hashedPassword),
		Name:     name,
	}
	_, err = service.UserRepository.Create(user)
	if err != nil {
		return 0, "", err
	}
	return user.ID, user.Email, err
}
