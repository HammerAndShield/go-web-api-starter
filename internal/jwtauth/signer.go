package jwtauth

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

// Signer holds the necessary variables to sign JWTs
type Signer struct {
	PrivateKey     *ecdsa.PrivateKey
	Issuer         string
	AcceptAudience []string
}

func parseECDSAPrivateKey(pemEncodedBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemEncodedBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse the PEM block containing the key")
	}

	privKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

// NewSigner creates a new instance of Signer using the provided base64 encoded string
// which is expected to be an ECDSA256 private key, issuer and an accepted audience slice.
// Returns an error if the provided base64 string is not a valid ECDSA256 private key.
// Also, returns an error if the curve used is not equal to elliptic.P256.
func NewSigner(base64String string, issuer string, acceptAudience []string) (*Signer, error) {
	key, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		return nil, err
	}

	pKey, err := parseECDSAPrivateKey(key)
	if err != nil {
		return nil, err
	}

	if pKey.Curve != elliptic.P256() {
		return nil, ErrWrongECDSAFormat
	}

	return &Signer{PrivateKey: pKey, Issuer: issuer, AcceptAudience: acceptAudience}, nil
}

// GenerateJWT creates a new JWT using the signers private key and returns it as a string.
// Will error if the key is not ECDSA256.
func (signer *Signer) GenerateJWT(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	tokenString, err := token.SignedString(signer.PrivateKey)
	if err != nil {
		return tokenString, err
	}

	return tokenString, nil
}
