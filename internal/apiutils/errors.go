package apiutils

import "errors"

var (
	ErrVersionNotInt        = errors.New("X-Expected-Version header must be an int")
	ErrVersionHeaderMissing = errors.New("X-Expected-Version header is missing")
)
