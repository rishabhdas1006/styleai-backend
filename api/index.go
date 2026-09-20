package handler

import (
	"fmt"
	"net/http"
	"os"

	"styleai-backend/pkg/server"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	setCORSHeaders(w, r)

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

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

func setCORSHeaders(w http.ResponseWriter, r *http.Request) {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		return
	}

	if r.Header.Get("Origin") == frontendURL {
		w.Header().Set("Access-Control-Allow-Origin", frontendURL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		w.Header().Set("Vary", "Origin")
	}
}
