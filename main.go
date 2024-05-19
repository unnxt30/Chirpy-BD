package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main(){
	var apiCfg apiConfig;
	server_mux := http.NewServeMux()
	var mock_server http.Server
	
	server_mux.Handle("/app/*", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("")))))

	
	server_mux.HandleFunc("GET /api/healthz", handleRequest);
	server_mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics);
	server_mux.HandleFunc("GET /api/reset", apiCfg.resetMetrics);
	server_mux.HandleFunc("POST /api/validate_chirp", validateChirp)


	mock_server.Addr = ":8080"
	mock_server.Handler = server_mux;

	mock_server.ListenAndServe();
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Set the status code to 404 Not Found
    w.WriteHeader(http.StatusOK);
	w.Header().Set("Content-type", "text/plain")
	w.Header().Add("Content-type", "charset=utf-8")
	fmt.Fprintf(w, "OK")	
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request){
	//fmt.Fprintf(w, "Hits: %v", cfg.fileserverHits);
	data, err := os.ReadFile("metrics.html")
		if err != nil {
			// Handle file reading error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the HTML file content as the response
		html_content := fmt.Sprintf(string(data), cfg.fileserverHits);
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html_content));
}

func(cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request){
	cfg.fileserverHits = 0;
}

type apiConfig struct {
	fileserverHits int;
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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {

    response, err := json.Marshal(payload)

    if err != nil {

        return err

    }

    w.WriteHeader(code)

    w.Write(response)

    return nil

}


func validateChirp(w http.ResponseWriter, r *http.Request){
	type parameters struct {
        Body string `json:"body"`
    }

	decoder := json.NewDecoder(r.Body);
	params := parameters{};
	err := decoder.Decode(&params);

	if err != nil {
		respondWithJSON(w, 404, map[string]string{"error": "Something went wrong"})
	}

	if len(params.Body) > 140 {
		respondWithJSON(w, 400, map[string]string{"error": "Chirp is too long"})
	}else{
		found := false
		var cleaned_body_words []string;
		bad_words := []string{"kerfuffle", "sharbert", "fornax"};
		//respondWithJSON(w, 200, map[string]bool{"valid" : true})
		input_str := params.Body;
		for _,word := range strings.Split(input_str, " "){
			found = false;
			for _, bad_word := range bad_words {
				if strings.ToLower(word) == bad_word{
					found = true;
					break;
				}
			}
			if found{
					cleaned_body_words = append(cleaned_body_words, "****");
			}else{
				cleaned_body_words = append(cleaned_body_words, word);
			}
		}	
		cleaned_body := strings.Join(cleaned_body_words, " ");
		respondWithJSON(w, 200, map[string]string{"cleaned_body": cleaned_body});
	}



}