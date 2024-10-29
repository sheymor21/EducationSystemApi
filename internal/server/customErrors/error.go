package customErrors

import (
	"SchoolManagerApi/internal/utilities"
	"fmt"
	"net/http"
)

type NotFoundMongoError struct {
	item string
	Msg  string
}

func NewNotFoundMongoError(item string) *NotFoundMongoError {
	err := &NotFoundMongoError{
		item: item,
		Msg:  fmt.Sprintf("Not Found this %s", item),
	}
	return err
}

func (err *NotFoundMongoError) Error() string {
	return fmt.Sprintf("Not Found this %s", err.item)
}

func ThrowHttpError(error error, w http.ResponseWriter, msg string, statusCode int) {
	if error == nil {
		return
	}

	if statusCode == http.StatusInternalServerError {
		utilities.Log.Errorln(error)
	}

	if msg != "" {
		http.Error(w, msg, statusCode)
	} else {
		http.Error(w, error.Error(), statusCode)
	}
	panic(nil)
}
