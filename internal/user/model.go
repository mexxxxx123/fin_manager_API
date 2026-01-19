package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string `gorm:"index"`
}

func NewUser(name string, password string, email string) *User {
	user := &User{
		Name:     name,
		Password: password,
		Email:    email,
	}
	return user
}
