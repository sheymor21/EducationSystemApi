package services

import (
	"net/http"
)

func httpInternalError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusInternalServerError)
}

func httpNotFoundError(w http.ResponseWriter, msg string) {
	http.Error(w, msg, http.StatusNotFound)
}
