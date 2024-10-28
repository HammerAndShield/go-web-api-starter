package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

// Signer holds the necessary variables to sign JWTs
type Signer struct {
	HmacSecret     []byte
	Issuer         string
	AcceptAudience []string
}

// NewSigner creates a new instance of Signer using the provided HMAC secret string,
// issuer and an accepted audience slice.
func NewSigner(hmacSecret string, issuer string, acceptAudience []string) (*Signer, error) {
	if hmacSecret == "" {
		return nil, errors.New("hmac secret cannot be empty")
	}

	return &Signer{
		HmacSecret:     []byte(hmacSecret),
		Issuer:         issuer,
		AcceptAudience: acceptAudience,
	}, nil
}

// GenerateJWT creates a new JWT using the signer's HMAC secret and returns it as a string.
// Uses HS256 (HMAC with SHA-256) for signing.
func (signer *Signer) GenerateJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signer.HmacSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
