package utilities

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"io"
	"log"
	"net/http"
)

type envelopeMap struct {
	Errors map[string]string
}

type envelopeMsg struct {
	Error string
}

func ReadJsonMiddlewareVersion(w http.ResponseWriter, r *http.Request, mapper any) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	var bodyCopy bytes.Buffer
	tee := io.TeeReader(r.Body, &bodyCopy)
	r.Body = io.NopCloser(&bodyCopy)
	enc := json.NewDecoder(tee)
	enc.DisallowUnknownFields()

	if err := enc.Decode(mapper); err != nil {
		return err
	}
	return nil
}

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

func WriteJsonError(w http.ResponseWriter, status int, msg interface{}) {

	var envelope any
	switch msg.(type) {
	case validator.ValidationErrors:
		envelope = envelopeMap{Errors: buildValidationMessages(msg.(validator.ValidationErrors))}
	case string:
		envelope = envelopeMsg{Error: msg.(string)}
	default:
		log.Fatal(errors.New("invalid type"))
	}
	indent, err := json.MarshalIndent(envelope, "", "\t")
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_, err = w.Write(indent)
	if err != nil {
		log.Fatal(err)
	}
}

func buildValidationMessages(errors validator.ValidationErrors) map[string]string {
	errorMap := make(map[string]string)

	for _, err := range errors {
		if err.ActualTag() == "min" {
			errorMap[err.Field()] = "minimum number of character" + " validation failed."
		} else if err.ActualTag() == "max" {

			errorMap[err.Field()] = "maximum number of character" + " validation failed."
		} else {
			errorMap[err.Field()] = err.ActualTag() + " validation failed."
		}
	}
	return errorMap
}
