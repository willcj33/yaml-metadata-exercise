package handlers

import (
	"net/http"
)

// Health Health check handler
func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/text")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ok"))
}
