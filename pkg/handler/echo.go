package handler

import (
	"encoding/json"
	"net/http"
)

type echoResponse struct {
	Message string `json:"message"`
}

// Echo returns the msg query parameter as JSON: {"message":"..."}.
// Example: GET /v1/echo?msg=hello
func Echo(w http.ResponseWriter, r *http.Request) {
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		http.Error(w, "missing 'msg' query parameter", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(echoResponse{Message: msg})
}
