package lib

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserPayload struct {
	Id   int    `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GeneratedTokens(id int, role string) string {
	secretKey := []byte(os.Getenv("APP_SECRET"))

	claims := UserPayload{
		Id:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(secretKey)
	return signedToken
}

func VerifyToken(token string) (UserPayload, error) {
	secretKey := []byte(os.Getenv("APP_SECRET"))
	parsedToken, err := jwt.ParseWithClaims(token, &UserPayload{}, func(t *jwt.Token) (any, error) {
		return secretKey, nil
	})
	if err != nil {
		return UserPayload{}, err
	} else if claims, ok := parsedToken.Claims.(*UserPayload); ok {
		return *claims, nil
	} else {
		return UserPayload{}, jwt.ErrTokenInvalidClaims
	}
}
