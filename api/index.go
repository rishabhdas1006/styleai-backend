package handler

import (
	"fmt"
	"net/http"

	"styleai-backend/pkg/server"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/health" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"StyleAI backend running"}`))
		return
	}

	app, err := server.Handler()
	if err != nil {
		http.Error(w, fmt.Sprintf("server initialization failed: %v", err), http.StatusInternalServerError)
		return
	}

	app.ServeHTTP(w, r)
}
