package server

import (
	"SchoolManagerApi/internal/services"
	"SchoolManagerApi/internal/utilities"
	"fmt"
	"github.com/MarceloPetrucio/go-scalar-api-reference"
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
	marksHandler := http.HandlerFunc(services.HttpMarksHandler)
	loginHandler := http.HandlerFunc(services.Login)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, scalarErr := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: app.swaggerSpecURL,
			CustomOptions: scalar.CustomOptions{
				PageTitle: "SchoolManager",
			},
			HideDownloadButton: true,
			DarkMode:           true,
		})
		if scalarErr != nil {
			utilities.Log.Warnln(scalarErr.Error())
		}

		_, printErr := fmt.Fprintln(w, htmlContent)
		if printErr != nil {
			utilities.Log.Warnln(printErr)
			return
		}
	})

	studentHandlerChain := loginMiddleware(loggerMiddleware(middlewareStudentValidator(studentHandler, app.validator)))
	studentsHandlerChain := loginMiddleware(loggerMiddleware(loggerMiddleware(loggerMiddleware(studentsHandler))))
	teacherHandlerChain := loginMiddleware(loggerMiddleware(middlewareTeacherValidator(teacherHandler, app.validator)))
	teachersHandlerChain := loginMiddleware(loginMiddleware(loggerMiddleware(teachersHandler)))
	markHandlerChain := loginMiddleware(loggerMiddleware(middlewareMarkValidator(markHandler, app.validator)))
	marksHandlerChain := loginMiddleware(loggerMiddleware(loggerMiddleware(marksHandler)))

	mux.Handle("/student", studentHandlerChain)
	mux.Handle("/students", studentsHandlerChain)
	mux.Handle("/teacher", teacherHandlerChain)
	mux.Handle("/teachers", teachersHandlerChain)
	mux.Handle("/mark", markHandlerChain)
	mux.Handle("/marks", marksHandlerChain)
	mux.Handle("/login", loginHandler)
	return mux
}
