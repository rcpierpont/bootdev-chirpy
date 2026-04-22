package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerHits(w http.ResponseWriter, r *http.Request) {
	hits := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(hits))
}
