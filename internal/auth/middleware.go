package auth

import (
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
            _, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
                return secret, nil
            })
            if err != nil {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
