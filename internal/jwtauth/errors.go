package jwtauth

import (
	"errors"
	"fmt"
)

var (
	ErrPemParse             = errors.New("error decoding pem from string")
	ErrNotECDSA             = errors.New("key is not ECDSA")
	ErrWrongECDSAFormat     = errors.New("invalid key: expected ES256 key")
	ErrImproperClaimsFormat = errors.New("claims are not formatted properly")
	ErrExpiredToken         = errors.New("provided tokens has expired")
	ErrInvalidIssuer        = errors.New("unaccepted issuer on tokens")
	ErrInvalidAudience      = errors.New("invalid tokens audience")
	ErrEmptySubject         = errors.New("subject of tokens is empty")
)

type ErrUnexpectedSigningMethod struct {
	unexpectedValue interface{}
}

func (e ErrUnexpectedSigningMethod) Error() string {
	return fmt.Sprintf("unexpected signing method %v", e.unexpectedValue)
}
