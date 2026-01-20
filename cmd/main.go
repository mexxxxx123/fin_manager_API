package main

import (
	"fin_manager_API/m/configs"
	"fin_manager_API/m/internal/auth"
	"fin_manager_API/m/internal/user"
	"fin_manager_API/m/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	mux := http.NewServeMux()

	//REPOSITORIES
	userRepo := user.NewUserRepository(db)

	//SERVICES
	authService := auth.NewAuthService(userRepo)

	//HANDLERS
	auth.NewAuthHandler(mux, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})

	server := http.Server{
		Addr:    ":8083",
		Handler: mux,
	}

	fmt.Println("леригоу ищу по порту 8083")

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err)
	}

}
