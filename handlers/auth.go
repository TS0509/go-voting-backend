package handlers

import "net/http"

func AuthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"ok"}`))
}

// AuthCheckHandler is a simple handler to check if the user is authenticated.
