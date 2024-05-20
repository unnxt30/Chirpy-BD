package src 

import (
	"fmt"
	"net/http"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Set the status code to 404 Not Found
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "text/plain")
	w.Header().Add("Content-type", "charset=utf-8")
	fmt.Fprintf(w, "OK")
}
