package validations

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	TeacherRol = iota + 1
	StudentRol
	AdminRol
)

type Rol int

type JWTUser struct {
	Carnet string
	Rol    Rol
}

func (r Rol) String() string {
	return [...]string{"Teacher", "Student", "Admin"}[r-1]
}

func (r Rol) EnumIndex() int {
	return int(r)
}

var secretKey []byte

func SetSecretKey(key []byte) {
	secretKey = key
}

func CreateJWT(user JWTUser) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Carnet": user.Carnet,
			"Rol":    user.Rol.String(),
			"exp":    time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, permissionRol []Rol) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if permissionRol != nil {
		permission := verifyRol(token, permissionRol)
		if !permission {
			return errors.New("permission Denied")
		}

	}

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func verifyRol(token *jwt.Token, roles []Rol) bool {

	for _, rol := range roles {
		if rol.String() == token.Claims.(jwt.MapClaims)["Rol"] {
			return true
		}
	}
	return false
}
