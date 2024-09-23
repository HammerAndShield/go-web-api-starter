package apiutils

import (
	"fmt"
	"net/http"
	"time"
)

func formatDuration(d time.Duration) string {
	return fmt.Sprintf("%.0f", d.Seconds())
}

func joinStrings(s []string) string {
	result := ""
	for i, v := range s {
		if i > 0 {
			result += ", "
		}
		result += v
	}
	return result
}

func GetCachingHeaders(maxAge time.Duration, staleWhileRevalidate time.Duration) http.Header {
	headers := make(http.Header)

	// Set Cache-Control header
	cacheControl := []string{
		"public",
		"max-age=" + formatDuration(maxAge),
		"stale-while-revalidate=" + formatDuration(staleWhileRevalidate),
	}
	headers.Set("Cache-Control", joinStrings(cacheControl))

	// Set Expires header
	expires := time.Now().Add(maxAge).Format(http.TimeFormat)
	headers.Set("Expires", expires)

	return headers
}
