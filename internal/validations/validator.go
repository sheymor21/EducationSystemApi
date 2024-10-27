package validations

import (
	"SchoolManagerApi/internal/utilities"
	"errors"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func Validate[T any](w http.ResponseWriter, r *http.Request, validate *validator.Validate, mapper T) error {

	err := utilities.ReadJsonMiddlewareVersion(w, r, &mapper)
	if err != nil {
		utilities.WriteJsonError(w, http.StatusBadRequest, err.Error())
		return err
	}
	structErr := validate.Struct(mapper)
	if structErr != nil {
		utilities.WriteJsonError(w, http.StatusBadRequest, structErr.(validator.ValidationErrors))
		return structErr
	}
	return nil
}

func LoginValidator(r *http.Request) (error error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return errors.New("missing authorization header")
	}
	tokenString = tokenString[len("Bearer "):]
	err := VerifyToken(tokenString)
	if err != nil {
		return errors.New("invalid Token")
	}
	return nil
}
