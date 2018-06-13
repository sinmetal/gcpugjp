package backend

import (
	"net/http"
)

func init() {
	http.HandleFunc("/_ah/warmup", handlerWarmup)
}

func handlerWarmup(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}