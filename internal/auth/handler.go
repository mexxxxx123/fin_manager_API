package auth

import "net/http"

type AuthHandler struct{}

func NewAuthHandler(mux *http.ServeMux) {
	handler := &AuthHandler{}
	mux.HandleFunc("/auth/login", handler.Login())
	mux.HandleFunc("/auth/register", handler.Register())
}

func (handler AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Login"))
	}

}
func (handler AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
