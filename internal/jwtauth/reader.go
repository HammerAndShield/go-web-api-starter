package jwtauth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Reader holds the necessary components to validate JWT tokens.
type Reader struct {
	HmacSecret []byte
	Issuer     string
}

// NewReader constructs a new Reader with the provided HMAC secret and issuer.
// These values are used when reading and validating JWT tokens.
// Requires the HMAC secret and issuer as string parameters.
// Returns a pointer to the new Reader and an error. If the construction is successful, the error is nil.
func NewReader(hmacSecret string, issuer string) (*Reader, error) {
	return &Reader{HmacSecret: []byte(hmacSecret), Issuer: issuer}, nil
}

// Read accepts a JWT token string, parses it, validates the signing method,
// and extracts its Claims as jwt.MapClaims. If the token's signing method is
// different from jwt.SigningMethodHMAC, ErrUnexpectedSigningMethod error is
// returned. If the claims inside the token cannot be formatted properly as
// jwt.MapClaims, ErrImproperClaimsFormat error is returned.
// It is recommended to validate the returned claims via reader's ValidateClaims method.
func (reader *Reader) Read(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrUnexpectedSigningMethod{unexpectedValue: token.Header["alg"]}
		}

		return reader.HmacSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrImproperClaimsFormat
	}

	return claims, nil
}

// ValidateClaims validates the provided JWT claims.
// It checks the expiration time, issuer, subject, and email fields.
// Returns appropriate errors if validation fails: ErrExpiredToken for expired tokens,
// ErrInvalidIssuer for an incorrect issuer, ErrEmptySubject for a missing subject,
// or a general error if the email claim is not a string.
func (reader *Reader) ValidateClaims(claims jwt.MapClaims) error {
	expiration, err := claims.GetExpirationTime()
	if err != nil {
		return fmt.Errorf("failed to get expiration time from claims: %w", err)
	}

	if time.Now().After(expiration.Time) {
		return ErrExpiredToken
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		return fmt.Errorf("failed to get issuer from claims: %w", err)
	}
	if issuer != reader.Issuer {
		return ErrInvalidIssuer
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return fmt.Errorf("failed to get subject from claims: %w", err)
	}
	if subject == "" {
		return ErrEmptySubject
	}

	_, ok := claims["email"].(string)
	if !ok {
		return errors.New("email claim is not a string")
	}

	return err
}
