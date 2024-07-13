package Utilities

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func ReadJson(body io.ReadCloser, mapper any) {
	all, err := io.ReadAll(body)
	if err != nil {
		return
	}

	err = json.Unmarshal(all, mapper)
	if err != nil {
		return
	}
}

func WriteJson(w http.ResponseWriter, status int, data any) {
	indent, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_, err = w.Write(indent)
	if err != nil {
		log.Fatal(err)
		return
	}
}
