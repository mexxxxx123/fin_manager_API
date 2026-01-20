package auth

import (
	"fin_manager_API/m/configs"
	"fin_manager_API/m/pkg/jwt"
	"fin_manager_API/m/pkg/req"
	"fin_manager_API/m/pkg/res"
	"fmt"
	"net/http"
)

type AuthHandler struct {
	*configs.Config
	*AuthService
}

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(mux *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	mux.HandleFunc("/auth/login", handler.Login())
	mux.HandleFunc("/auth/register", handler.Register())
}

func (handler AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			fmt.Println(err)
		}
		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		j := jwt.NewJWT(handler.Config.Auth.Secret)
		token, err := j.Create(jwt.JWTData{
			Email: email,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(LoginResponse{
			Token: token,
		}, w, http.StatusOK)

		fmt.Println("Login succeful!")
		fmt.Println(email, err)
	}

}
func (handler AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		j := jwt.NewJWT(handler.Config.Auth.Secret)
		token, err := j.Create(jwt.JWTData{
			Email: email,
		})
		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Json(RegisterResponse{
			Token: token,
		}, w, http.StatusOK)

		fmt.Println("Register succeful!")
		fmt.Println(email, err)
	}
}
