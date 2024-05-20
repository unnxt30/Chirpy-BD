package src 

import (
	"fmt"
	"net/http"
	"os"
)

type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	//fmt.Fprintf(w, "Hits: %v", cfg.fileserverHits);
	data, err := os.ReadFile("metrics.html")
	if err != nil {
		// Handle file reading error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send the HTML file content as the response
	html_content := fmt.Sprintf(string(data), cfg.fileserverHits)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html_content))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Increment the counter
		cfg.fileserverHits++

		// Log to see the count (for debugging, not required)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}
