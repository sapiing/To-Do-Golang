package middleware

import (
	"net/http"
	"strings"

	"backend/auth"
	"backend/error"
)

// Middleware untuk otentikasi
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            error.HandleError(w, http.StatusUnauthorized, "Authorization header missing", nil)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // Validasi token
        _, err := auth.ValidationToken(tokenString)
        if err != nil {
            error.HandleError(w, http.StatusUnauthorized, "Invalid token", err)
            return
        }

        // Token valid, lanjutkan ke handler berikutnya
        next.ServeHTTP(w, r)
    })
}
