package main

import (
	"fin_manager_API/m/configs"
	"fin_manager_API/m/internal/auth"
	"fin_manager_API/m/internal/stat"
	"fin_manager_API/m/internal/transactions"
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
	transRepo := transactions.NewTransactionRepository(db)
	statRepo := stat.NewStatRepository(db)

	//SERVICES
	authService := auth.NewAuthService(userRepo)
	transService := transactions.NewTrasactionService(transRepo)

	//HANDLERS
	transactions.NewTransactionHandler(mux, transactions.TransactionHandlerDeps{
		Config:                conf,
		TransactionRepository: transRepo,
		TransactionService:    transService,
	})
	auth.NewAuthHandler(mux, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	stat.NewStatHandler(mux, stat.StatHandlerDeps{
		Config:   conf,
		StatRepo: statRepo,
	})

	server := http.Server{
		Addr:    ":8084",
		Handler: mux,
	}

	fmt.Println("леригоу ищу по порту 8084")

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err)
	}

}
