package main

import (
	"fin_manager_API/m/internal/transactions"
	"fin_manager_API/m/internal/user"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	db.AutoMigrate(&user.User{}, &transactions.Transaction{})
}
