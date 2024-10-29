package server

import (
	"SchoolManagerApi/internal/dto"
	"SchoolManagerApi/internal/utilities"
	"SchoolManagerApi/internal/validations"
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
			teacherDto := dto.TeacherAddRequest{}
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
			studentAddDto := dto.StudentAddRequest{}
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
		}).Infoln("Completed request")

	})
}

func loginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer utilities.Recover()
		switch r.URL.Path {
		case "/students":
			if r.Method == http.MethodGet {
				validatePermissions(w, r, validations.TeacherRol, validations.AdminRol)
			}
		case "/Student":
			validatePermissions(w, r, validations.TeacherRol, validations.AdminRol)
		case "/marks":
			validatePermissions(w, r, validations.TeacherRol, validations.AdminRol)
		case "/mark":
			validatePermissions(w, r, validations.TeacherRol, validations.AdminRol)
		case "/teacher":
			if r.Method != http.MethodPost {
				validatePermissions(w, r, validations.TeacherRol, validations.AdminRol)
			}
		case "/teachers":
			validatePermissions(w, r, validations.TeacherRol, validations.AdminRol)
		default:
			validatePermissions(w, r)

		}
		next.ServeHTTP(w, r)
	})

}

func validatePermissions(w http.ResponseWriter, r *http.Request, permissions ...validations.Rol) {
	validationErr := validations.LoginValidator(r, permissions)
	if validationErr != nil {
		utilities.WriteJsonError(w, http.StatusUnauthorized, validationErr.Error())
		panic("")
	}
}
