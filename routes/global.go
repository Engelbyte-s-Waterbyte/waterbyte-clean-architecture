package routes

import (
	"encoding/json"
	"net/http"
)

func badRequest(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	errResponse(w, "Bad Request")
}

func errResponse(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(map[string]interface{}{"error": message})
}

func jsonResponse(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}
