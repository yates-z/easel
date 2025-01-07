package jwt

import "errors"

var (
	ErrTokenHeaderInvalid  = errors.New("token header is invalid")
	ErrTokenPayloadInvalid = errors.New("token payload is invalid")
	ErrTokenInvalid        = errors.New("token is invalid")
	ErrMethodNotFound      = errors.New("signing method not found")
	ErrHashUnavailable     = errors.New("the requested hash function is unavailable")
	ErrSignatureInvalid    = errors.New("signature is invalid")
)
