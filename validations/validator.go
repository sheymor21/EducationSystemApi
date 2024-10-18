package validations

import (
	"calificationApi/internal/utilities"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func Validate(w http.ResponseWriter, r *http.Request, validate *validator.Validate, mapper any) error {
	err := utilities.ReadJsonMiddlewareVersion(w, r, &mapper)
	if err != nil {
		return err
	}
	err = validate.Struct(mapper)
	if err != nil {
		utilities.WriteJsonError(w, http.StatusBadRequest, err.(validator.ValidationErrors))
		return err
	}
	return nil
}
