package di

import "fin_manager_API/m/internal/user"

type IUserRepository interface {
	FindByEmail(email string) (*user.User, error)
	Create(*user.User) (*user.User, error)
}
