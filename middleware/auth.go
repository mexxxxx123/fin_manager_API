package middleware

import (
	"context"
	"errors"
	"fin_manager_API/m/configs"
	"fin_manager_API/m/pkg/jwt"
	"fmt"
	"net/http"
	"strings"
)

type key string

const (
	ContextUserIdKey key = "ContextUserIdKey"
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
		ctx := context.WithValue(r.Context(), ContextUserIdKey, data.UserId)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

func GetUserID(ctx context.Context) (uint, error) {
	userID, ok := ctx.Value(ContextUserIdKey).(uint)
	if !ok {
		return 0, errors.New("cant get userId from context")
	}
	return userID, nil
}

// func IsAuthed(next http.Handler, conf *configs.Config) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("=== MIDDLEWARE DEBUG ===")
// 		fmt.Printf("URL: %s %s\n", r.Method, r.URL.Path)
//
// 		autorizationStr := r.Header.Get("Authorization")
// 		fmt.Printf("Authorization header: %q\n", autorizationStr)
//
// 		if !strings.HasPrefix(autorizationStr, "Bearer ") {
// 			fmt.Println("❌ No Bearer prefix")
// 			writeUnauthed(w)
// 			return
// 		}
//
// 		if autorizationStr == "" {
// 			fmt.Println("❌ Empty Authorization")
// 			writeUnauthed(w)
// 			return
// 		}
//
// 		_, token, found := strings.Cut(autorizationStr, "Bearer ")
// 		fmt.Printf("Found token: %v, token: %s...\n", found, token[:min(30, len(token))])
//
// 		isValid, data := jwt.NewJWT(conf.Auth.Secret).Parse(token)
// 		fmt.Printf("Parse result: isValid=%v, data=%+v\n", isValid, data)
//
// 		if !isValid {
// 			fmt.Println("❌ Token not valid in middleware")
// 			writeUnauthed(w)
// 			return
// 		}
//
// 		fmt.Printf("✅ User authenticated: ID=%d, Email=%s\n", data.UserId, data.Email)
// 		ctx := context.WithValue(r.Context(), ContextUserIdKey, data.UserId)
// 		req := r.WithContext(ctx)
// 		next.ServeHTTP(w, req)
// 	})
// }
