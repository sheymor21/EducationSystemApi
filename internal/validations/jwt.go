package validations

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type JWTUser struct {
	UserName string
	Password string
}

var secretKey []byte

func SetSecretKey(key []byte) {
	secretKey = key
}

func CreateJWT(user JWTUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"UserName": user.UserName,
			"Password": user.Password,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}
