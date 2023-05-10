package auth

import (
	"fmt"
)

const (
	ErrCode = iota + 10000
	ErrCodeBadIssuer
	ErrCodeBadDomain
	ErrCodeBadLoginMethod
	ErrCodeBadLoginMessage
	ErrCodeBadLoginSignature
)

type (
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message,omitempty"`
	}
)

func (err *Error) Error() string {
	return fmt.Sprintf("[%d] %s", err.Code, err.Message)
}

func NewError(msg string) error {
	return &Error{
		Code:    ErrCode,
		Message: msg,
	}
}

func IsError(err error, code int) bool {
	if err, ok := err.(*Error); ok {
		return err.Code == code
	}
	return false
}

func NewBadIssuerError(hint string) error {
	msg := "bad issuer"
	if hint != "" {
		msg = fmt.Sprintf("%s: %s", msg, hint)
	}
	return &Error{
		Code:    ErrCodeBadIssuer,
		Message: msg,
	}
}

func IsErrBadIssuer(err error) bool {
	return IsError(err, ErrCodeBadIssuer)
}

func NewBadDomainError(hint string) error {
	msg := "bad domain"
	if hint != "" {
		msg = fmt.Sprintf("%s: %s", msg, hint)
	}
	return &Error{
		Code:    ErrCodeBadDomain,
		Message: msg,
	}
}

func IsErrBadDomain(err error) bool {
	return IsError(err, ErrCodeBadDomain)
}

func NewBadLoginMethodError(hint string) error {
	msg := "bad login method"
	if hint != "" {
		msg = fmt.Sprintf("%s: %s", msg, hint)
	}
	return &Error{
		Code:    ErrCodeBadLoginMethod,
		Message: msg,
	}
}

func IsErrBadLoginMethod(err error) bool {
	return IsError(err, ErrCodeBadLoginMethod)
}

func NewBadLoginMessageError(hint string) error {
	msg := "bad mvm message"
	if hint != "" {
		msg = fmt.Sprintf("%s: %s", msg, hint)
	}
	return &Error{
		Code:    ErrCodeBadLoginMessage,
		Message: msg,
	}
}

func IsErrBadLoginMessage(err error) bool {
	return IsError(err, ErrCodeBadLoginMessage)
}

func NewBadLoginSignatureError(hint string) error {
	msg := "bad mvm signature"
	if hint != "" {
		msg = fmt.Sprintf("%s: %s", msg, hint)
	}
	return &Error{
		Code:    ErrCodeBadLoginSignature,
		Message: msg,
	}
}

func IsErrBadLoginSignature(err error) bool {
	return IsError(err, ErrCodeBadLoginSignature)
}
