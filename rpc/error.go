package rpc

import (
	"errors"
)

type ErrorCode int

const (
	ErrorCodeParse            ErrorCode = -32700
	ErrorCodeInvalidIDRequest ErrorCode = -32600
	ErrorCodeNoMethod         ErrorCode = -32601
	ErrorCodeBadParams        ErrorCode = -32602
	ErrorCodeInternal         ErrorCode = -32603
	ErrorCodeServer           ErrorCode = -32000
)

var ErrorNullResult = errors.New("result is null")

type Error struct {
	Code    ErrorCode   `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (e *Error) Error() string {
	return e.Message
}
