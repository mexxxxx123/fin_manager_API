package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWT struct {
	Secret string
}

type JWTData struct {
	Email  string
	UserId uint
}

func NewJWT(secret string) *JWT {
	return &JWT{
		Secret: secret,
	}
}
func (j *JWT) Create(data JWTData) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":  data.Email,
		"userId": data.UserId,
	})
	s, err := t.SignedString([]byte(j.Secret))
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) Parse(token string) (bool, *JWTData) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return []byte(j.Secret), nil
	})
	if err != nil {
		fmt.Println("❌ Parse failed:", err)
		return false, nil
	}
	if !t.Valid {
		fmt.Println("❌ Token not valid")
		return false, nil
	}
	email := t.Claims.(jwt.MapClaims)["email"]
	userId := t.Claims.(jwt.MapClaims)["userId"]
	userIdFloat := userId.(float64)
	return t.Valid, &JWTData{
		Email:  email.(string),
		UserId: uint(userIdFloat),
	}

}

// func (j *JWT) Parse(token string) (bool, *JWTData) {
// 	fmt.Println("=== JWT DEBUG ===")
// 	fmt.Printf("Token length: %d\n", len(token))
// 	fmt.Printf("Token preview: %s...\n", token[:min(30, len(token))])
//
// 	t, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
// 		return []byte(j.Secret), nil
// 	})
//
// 	fmt.Printf("Parse error: %v\n", err)
// 	fmt.Printf("Token valid: %v\n", t != nil && t.Valid)
//
// 	if err != nil {
// 		fmt.Println("❌ Parse failed:", err)
// 		return false, nil
// 	}
//
// 	if !t.Valid {
// 		fmt.Println("❌ Token not valid")
// 		return false, nil
// 	}
//
// 	claims := t.Claims.(jwt.MapClaims)
// 	fmt.Printf("Claims: %+v\n", claims)
//
// 	emailClaim, emailOk := claims["email"]
// 	userIdClaim, userIdOk := claims["userId"]
//
// 	fmt.Printf("Email claim exists: %v, type: %T, value: %v\n", emailOk, emailClaim, emailClaim)
// 	fmt.Printf("UserId claim exists: %v, type: %T, value: %v\n", userIdOk, userIdClaim, userIdClaim)
//
// 	if !emailOk || !userIdOk || emailClaim == nil || userIdClaim == nil {
// 		fmt.Println("❌ Claims missing or nil")
// 		return false, nil
// 	}
//
// 	email, emailOk := emailClaim.(string)
// 	if !emailOk {
// 		fmt.Printf("❌ Email not string: %T\n", emailClaim)
// 		return false, nil
// 	}
//
// 	userIdFloat, userIdOk := userIdClaim.(float64)
// 	if !userIdOk {
// 		fmt.Printf("❌ UserId not float64: %T\n", userIdClaim)
// 		return false, nil
// 	}
//
// 	fmt.Printf("✅ Parsed: email=%s, userId=%v (converted to uint=%v)\n",
// 		email, userIdFloat, uint(userIdFloat))
//
// 	return true, &JWTData{
// 		Email:  email,
// 		UserId: uint(userIdFloat),
// 	}
// }
//
