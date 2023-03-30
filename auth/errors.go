package auth

import "errors"

var (
	ErrInvalidIssuer        = errors.New("invalid issuer")
	ErrBadLoginMethod       = errors.New("bad login method")
	ErrBadMvmLoginSignature = errors.New("bad mvm login signature")
	ErrBadMvmLoginMessage   = errors.New("bad mvm login message")
)
