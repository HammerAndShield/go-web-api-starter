package middleware

import "net/http"

// Chain composes middlewares into a single handler.
func Chain(final http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	if len(middlewares) == 0 {
		return final
	}
	wrapped := final
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i](wrapped)
	}
	return wrapped
}
