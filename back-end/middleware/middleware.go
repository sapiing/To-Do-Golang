package middleware

import (
    "net/http"

    "backend/auth"
)

// Middleware untuk otentikasi
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // ambil token dari header
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
            return
        }

        // pisahin token dari prefix bearer
        tokenString := authHeader[len("Bearer "):]

        // validasi token
        _, err := auth.ValidationToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // token valid, lanjut ke header selanjutnya
        next.ServeHTTP(w, r)
    })
}
