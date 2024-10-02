package server

import (
	"calificationApi/internal/services"
	"net/http"
)

func (app *application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	studentHandler := http.HandlerFunc(services.HttpStudentHandler)
	teacherHandler := http.HandlerFunc(services.HttpTeacherHandler)
	markHandler := http.HandlerFunc(services.HttpMarkHandler)

	mux.Handle("/student", middlewareStudentValidator(studentHandler, app.validator))
	mux.Handle("/teacher", middlewareTeacherValidator(teacherHandler, app.validator))
	mux.Handle("/mark", middlewareMarkValidator(markHandler, app.validator))
	mux.HandleFunc("/marks", services.HttpMarksHandler)
	return mux
}
