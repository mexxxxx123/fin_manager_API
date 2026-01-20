package middleware

import (
	"context"
	"fin_manager_API/m/configs"
	"fin_manager_API/m/pkg/jwt"
	"fmt"
	"net/http"
	"strings"
)

type key string

const (
	ContextEmailKey key = "ContextEmailKey"
)

func writeUnauthed(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte(http.StatusText(http.StatusUnauthorized)))
}

func IsAuthed(next http.Handler, conf *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		autorizationStr := r.Header.Get("Authorization")
		if !strings.HasPrefix(autorizationStr, "Bearer ") {
			writeUnauthed(w)
			return
		}

		if autorizationStr == "" {
			fmt.Println("BEARERTOKEN_NOTFOUND")
			next.ServeHTTP(w, r)
			return
		}
		_, token, _ := strings.Cut(autorizationStr, "Bearer ")
		isValid, data := jwt.NewJWT(conf.Auth.Secret).Parse(token)
		if !isValid {
			writeUnauthed(w)
			return
		}
		ctx := context.WithValue(r.Context(), ContextEmailKey, data.Email)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
