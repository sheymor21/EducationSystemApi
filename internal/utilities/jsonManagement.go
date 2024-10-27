package utilities

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io"
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
		Log.Errorln(err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_, err = w.Write(indent)
	if err != nil {
		Log.Errorln(err)
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
		Log.Errorln(errors.New("invalid type"))
	}
	indent, err := json.MarshalIndent(envelope, "", "\t")
	if err != nil {
		Log.Fatalln(err)
	}
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	w.WriteHeader(status)
	_, err = w.Write(indent)
	if err != nil {
		Log.Fatalln(err)
	}
}

func buildValidationMessages(errors validator.ValidationErrors) map[string]string {
	errorMap := make(map[string]string)

	for _, err := range errors {
		errorMap[getFieldValidationName(err)] = getFieldValidationTag(err)
	}
	return errorMap
}
func getFieldValidationTag(err validator.FieldError) string {
	switch err.Tag() {
	case "min":
		return fmt.Sprintf("the minimum number of character is %s", err.Param())
	case "max":
		return fmt.Sprintf("the maximum number of character is %s", err.Param())
	default:
		return err.ActualTag() + " validation failed"
	}
}
func getFieldValidationName(err validator.FieldError) string {
	switch err.Field() {
	case "StudentCarnet":
		return "Student_Carnet"
	case "TeacherCarnet":
		return "Teacher_Carnet"
	default:
		return err.Field()
	}
}
