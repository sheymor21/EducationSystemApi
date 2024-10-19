package server

import (
	"calificationApi/internal/dto"
	"calificationApi/internal/utilities"
	"calificationApi/validations"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

func middlewareMarkValidator(next http.Handler, validate *validator.Validate) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			markDto := dto.MarkAddRequest{}
			err := validations.Validate(w, r, validate, markDto)
			if err != nil {
				return
			}

		case http.MethodPut:
			markDto := dto.MarksUpdateRequest{}
			err := validations.Validate(w, r, validate, markDto)
			if err != nil {
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
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		start := time.Now()
		next.ServeHTTP(w, r)
		utilities.Log.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"duration": time.Since(start),
		}).Info("Completed request")

	})
}
