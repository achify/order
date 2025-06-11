package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Middleware validates JWT token from Authorization header
func Middleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			rolesIface, _ := claims["roles"].([]interface{})
			roles := make([]string, 0, len(rolesIface))
			for _, v := range rolesIface {
				if s, ok := v.(string); ok {
					roles = append(roles, s)
				}
			}
			sub, _ := claims["sub"].(string)
			ctx := context.WithValue(r.Context(), ctxUserID{}, sub)
			ctx = context.WithValue(ctx, ctxRoles{}, roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type ctxUserID struct{}
type ctxRoles struct{}

// UserIDFromContext returns authenticated user id from context.
func UserIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(ctxUserID{}).(string); ok {
		return v
	}
	return ""
}

// RolesFromContext returns role list from context.
func RolesFromContext(ctx context.Context) []string {
	if v, ok := ctx.Value(ctxRoles{}).([]string); ok {
		return v
	}
	return nil
}
