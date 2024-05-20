package main

import (
	"net/http"

	"github.com/unnxt30/Chirpy-BD/src"
)


func main(){
	var apiCfg src;
	server_mux := http.NewServeMux()
	var mock_server http.Server
	
	server_mux.Handle("/app/*", http.StripPrefix("/app/", apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("")))))

	
	server_mux.HandleFunc("GET /api/healthz", handleRequest);
	server_mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics);
	server_mux.HandleFunc("GET /api/reset", apiCfg.resetMetrics);
	server_mux.HandleFunc("POST /api/chirps", validateChirp)
	//server_mux.HandleFunc("GET /api/chirps")

	mock_server.Addr = ":8080"
	mock_server.Handler = server_mux;

	mock_server.ListenAndServe();
}





