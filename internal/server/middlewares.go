package server

import (
	"calificationApi/internal/dto"
	"calificationApi/internal/utilities"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

func middlewareMarkValidator(next http.Handler, validate *validator.Validate) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			markDto := dto.MarkAddRequest{}
			err := utilities.ReadJsonMiddlewareVersion(w, r, &markDto)
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			err = validate.Struct(markDto)
			if err != nil {
				utilities.WriteJsonError(w, http.StatusBadRequest, err.(validator.ValidationErrors))
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
			err := utilities.ReadJsonMiddlewareVersion(w, r, &teacherDto)
			if err != nil {
				return
			}
			err = validate.Struct(teacherDto)
			if err != nil {
				utilities.WriteJson(w, http.StatusBadRequest, err.(validator.ValidationErrors))
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
			log.Println(r.URL.Path)
			studentAddDto := dto.StudentAddDto{}
			err := utilities.ReadJsonMiddlewareVersion(w, r, &studentAddDto)
			if err != nil {
				utilities.WriteJsonError(w, http.StatusBadRequest, err.Error())
				return
			}
			err = validate.Struct(studentAddDto)
			if err != nil {
				utilities.WriteJsonError(w, http.StatusBadRequest, err.(validator.ValidationErrors))
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
