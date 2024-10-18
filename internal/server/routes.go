package server

import (
	"calificationApi/internal/services"
	"fmt"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"log"
	"net/http"
)

// Routes sets up the routes and handlers for the application and returns an *http.ServeMux configured with them.
func (app *application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	studentHandler := http.HandlerFunc(services.HttpStudentHandler)
	studentsHandler := http.HandlerFunc(services.HttpStudentsHandler)
	teacherHandler := http.HandlerFunc(services.HttpTeacherHandler)
	teachersHandler := http.HandlerFunc(services.HttpTeachersHandler)
	markHandler := http.HandlerFunc(services.HttpMarkHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, scalarErr := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: app.swaggerSpecURL,
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Simple API",
			},
			DarkMode: true,
		})

		if scalarErr != nil {
			fmt.Printf("%v", scalarErr)
		}

		_, printErr := fmt.Fprintln(w, htmlContent)
		if printErr != nil {
			log.Println(printErr)
			return
		}
	})

	mux.Handle("/student", middlewareStudentValidator(studentHandler, app.validator))
	mux.Handle("/students", studentsHandler)
	mux.Handle("/teacher", middlewareTeacherValidator(teacherHandler, app.validator))
	mux.Handle("/teachers", teachersHandler)
	mux.Handle("/mark", middlewareMarkValidator(markHandler, app.validator))
	mux.HandleFunc("/marks", services.HttpMarksHandler)
	return mux
}
