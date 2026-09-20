package handler

import (
	"fmt"
	"net/http"

	"styleai-backend/pkg/server"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	app, err := server.Handler()
	if err != nil {
		http.Error(w, fmt.Sprintf("server initialization failed: %v", err), http.StatusInternalServerError)
		return
	}

	app.ServeHTTP(w, r)
}
