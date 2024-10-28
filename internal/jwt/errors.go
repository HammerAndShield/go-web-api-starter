package jwt

import (
	"errors"
	"fmt"
)

var (
	ErrImproperClaimsFormat = errors.New("claims are not formatted properly")
	ErrExpiredToken         = errors.New("provided tokens has expired")
	ErrInvalidIssuer        = errors.New("unaccepted issuer on tokens")
	ErrEmptySubject         = errors.New("subject of tokens is empty")
)

type ErrUnexpectedSigningMethod struct {
	unexpectedValue interface{}
}

func (e ErrUnexpectedSigningMethod) Error() string {
	return fmt.Sprintf("unexpected signing method %v", e.unexpectedValue)
}
