package error

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleError(w http.ResponseWriter, statusCode int, message string, err error) {
	log.Println(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
