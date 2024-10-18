package server

import (
	"calificationApi/internal/dto"
	"calificationApi/internal/utilities"
	"calificationApi/validations"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func middlewareMarkValidator(next http.Handler, validate *validator.Validate) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			markDto := dto.MarkAddRequest{}

			err := validations.Validate(w, r, validate, markDto)
			if err != nil {
				utilities.WriteJsonError(w, http.StatusBadRequest, err.Error())
				return
			}
		case http.MethodPut:
			markDto := dto.MarksUpdateRequest{}
			err := validations.Validate(w, r, validate, markDto)
			if err != nil {
				utilities.WriteJsonError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func middlewareTeacherValidator(next http.Handler, validate *validator.Validate) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			teacherDto := dto.TeacherDto{}
			err := validations.Validate(w, r, validate, teacherDto)
			if err != nil {
				utilities.WriteJsonError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func middlewareStudentValidator(next http.Handler, validate *validator.Validate) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			studentAddDto := dto.StudentAddDto{}
			err := validations.Validate(w, r, validate, studentAddDto)
			if err != nil {
				utilities.WriteJsonError(w, http.StatusBadRequest, err.Error())
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
