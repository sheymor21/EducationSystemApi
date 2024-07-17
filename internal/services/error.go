package services

import (
	"fmt"
)

type NotFoundMongoError struct {
	item string
	msg  string
}

func NewNotFoundMongoError(item string) *NotFoundMongoError {
	var err NotFoundMongoError
	err.msg = fmt.Sprintf("Not Found this %s", item)
	return &err
}

func (err *NotFoundMongoError) Error() string {
	return fmt.Sprintf("Not Found this %s", err.item)
}
