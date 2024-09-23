package apiutils

import (
	"net/http"
	"strconv"
)

const (
	expectedVersionHeader = "X-Expected-Version"
)

func CheckExpectedVersion(r *http.Request) (int, error) {
	expectedVersionH := r.Header.Get("X-Expected-Version")
	if expectedVersionH == "" {
		return 0, ErrVersionHeaderMissing
	}

	eVersion, err := strconv.Atoi(expectedVersionH)
	if err != nil {
		return 0, ErrVersionNotInt
	}

	return eVersion, nil
}
