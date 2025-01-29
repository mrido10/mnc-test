package model

import (
	"errors"
	"github.com/mrido10/ido-log/runtime"
)

type Error struct {
	error
	Code      int
	Message   string
	ErrorFile string
}

func NewError(code int, msg string, err error) *Error {
	if err == nil {
		err = errors.New(msg)
	}
	return &Error{
		error:     err,
		Code:      code,
		Message:   msg,
		ErrorFile: "\n\t" + runtime.GetRuntimeCaller(2),
	}
}
