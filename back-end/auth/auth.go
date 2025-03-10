package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	username  = "admin"
	password  = "admin"
	secretKey = "rahasia"
)

func GenerateToken() (string, error) {
	// ngambil token
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// bikin token baru
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// secret token
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidationToken(signedToken string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// validasi tambahan
	if claims["username"] != username {
		return nil, fmt.Errorf("invalid username")
	}

	// return tokennya
	return claims, nil
}

func HandleToken(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// basic authentication
	if creds.Username != username || creds.Password != password {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid credentials"})
		return
	}

	tokenString, err := GenerateToken()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to generate token"})
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

func HandleLogout(w http.ResponseWriter, r *http.Request) {

    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")  // For development
    w.Header().Set("Access-Control-Allow-Credentials", "true")
    // Clear the token from the client-side (browser)
    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   "",
        Path:    "/",
        Expires: time.Unix(0, 0), // Expire immediately
        MaxAge:  -1,             // Delete the cookie
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
