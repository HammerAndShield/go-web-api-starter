package users

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type contextKey string

const userContextKey = contextKey("user")

// contextSetUser associates the provided *data.User with the *http.Request using Context.
// It can be later retrieved in other parts of the code that have access to this http.Request using contextGetUser.
// This operation does not modify the incoming http.Request but instead returns a new http.Request
// with the new Context. The original http.Request should be discarded and the returned http.Request should be used thereafter.
func contextSetUser(r *http.Request, user *User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// ContextGetUser retrieves the *data.User associated with the *http.Request.
// This method will panic with the message "missing user value in request context"
// if no user value is available in the request context. This could happen if the contextSetUser method
// was not called to associate a User value with this request, or if the value was associated
// but is not of the expected *data.User type.
func ContextGetUser(r *http.Request) *User {
	user, ok := r.Context().Value(userContextKey).(*User)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}

// ContextGetUserId functions similarly to ContextGetUser but only returns the ID.
// Panics if the user is not found in the context.
func ContextGetUserId(r *http.Request) uuid.UUID {
	user := ContextGetUser(r)
	return user.ID
}
