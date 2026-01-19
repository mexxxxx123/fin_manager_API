package main

import (
	"fin_manager_API/m/configs"
	"fin_manager_API/m/internal/auth"
	"fin_manager_API/m/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	_ = db.NewDb(conf)
	mux := http.NewServeMux()

	//HANDLERS
	auth.NewAuthHandler(mux)

	server := http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	fmt.Println("леригоу ищу по порту 8081")

	err := server.ListenAndServe()

	if err != nil {
		fmt.Println(err)
	}

}
