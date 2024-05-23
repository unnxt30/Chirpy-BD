package src

import "net/http"

func (cfg *ApiConfig) ResetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
}
