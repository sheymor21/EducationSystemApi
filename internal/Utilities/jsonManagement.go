package Utilities

import (
	"encoding/json"
	"log"
	"net/http"
)

func ReadJson(w http.ResponseWriter, r *http.Request, mapper any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	enc := json.NewDecoder(r.Body)
	enc.DisallowUnknownFields()

	if err := enc.Decode(mapper); err != nil {
		return err
	}
	return nil
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
