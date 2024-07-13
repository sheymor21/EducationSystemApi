package server

import (
	"calificationApi/internal/services"
	"net/http"
)

func (app *application) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/student", services.HttpStudentHandler)
	return mux
}
