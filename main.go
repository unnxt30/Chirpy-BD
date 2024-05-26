package main

import (
	"net/http"

	"github.com/unnxt30/Chirpy-BD/src"
)

func main() {

	var apiCfg src.ApiConfig
	server_mux := http.NewServeMux()
	var mock_server http.Server

	server_mux.Handle("/app/*", http.StripPrefix("/app/", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir("")))))

	server_mux.HandleFunc("GET /api/healthz", src.HandleRequest)
	server_mux.HandleFunc("GET /admin/metrics", apiCfg.HandleMetrics)
	server_mux.HandleFunc("GET /api/reset", apiCfg.ResetMetrics)
	server_mux.HandleFunc("POST /api/chirps", src.ValidateChirp)
	server_mux.HandleFunc("GET /api/chirps", src.ChirpsGET)
	server_mux.HandleFunc("GET /api/chirps/{id}", src.ChirpGETbyID)
	server_mux.HandleFunc("POST /api/users", src.WriteUser)
	server_mux.HandleFunc("POST /api/login", src.LoginUser)
	mock_server.Addr = ":8080"
	mock_server.Handler = server_mux

	mock_server.ListenAndServe()
}
