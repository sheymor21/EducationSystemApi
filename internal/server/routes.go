package server

import (
	"calificationApi/internal/services"
	"net/http"
)

func (app *application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/student", services.HttpStudentHandler)
	mux.HandleFunc("/teacher", services.HttpTeacherHandler)
	mux.HandleFunc("/mark", services.HttpMarkHandler)
	mux.HandleFunc("/marks", services.HttpMarksHandler)
	return mux
}
