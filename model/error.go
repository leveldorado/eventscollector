package model

import (
	"reflect"
	"runtime"
)

const (
	ERR_ID_INVALID_DATA          = "INVALID_DATA"
	ERR_ID_INTERNAL_SERVER_ERROR = "INTERNAL_SERVER_ERROR"
	ERR_ID_STORE_ERROR           = "STORE_ERROR"
)

type AppError struct {
	Code    int    `json:"code"`
	ID      string `json:"id"`
	Message string `json:"message"`
	Method  string `json:"method"`
	File    string `json:"file"`
	Line    int    `json:"line"`
}

func NewAppError(id, message string, code int, method interface{}) *AppError {
	err := &AppError{
		Code:    code,
		ID:      id,
		Message: message,
		Method:  runtime.FuncForPC(reflect.ValueOf(method).Pointer()).Name(),
	}
	_, err.File, err.Line, _ = runtime.Caller(1)
	return err
}
