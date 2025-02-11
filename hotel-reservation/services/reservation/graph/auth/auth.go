// https://gqlgen.com/recipes/authentication/
package auth

import (
	"context"
	"net/http"
	"strconv"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

// A stand-in for our database backed user object
type User struct {
	ID   int
	Name string
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if user, err := getUserByID(r.Header.Get("user-id")); err == nil {
				ctx := context.WithValue(r.Context(), userCtxKey, user)
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *User {
	raw, _ := ctx.Value(userCtxKey).(*User)
	return raw
}

func validateAndGetUserID(c *http.Cookie) (string, error) {
	return c.Value, nil
}

func getUserByID(id string) (*User, error) {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:   idInt,
		Name: "user",
	}, nil
}
