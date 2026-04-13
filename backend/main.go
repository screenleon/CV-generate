package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/screenleon/cv-generate/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Resolve frontend directory: allow override via env, default to ../frontend
	// relative to the directory containing this source file (works with go run).
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		// When run with `go run ./backend` from repo root or with `go run main.go`
		// from the backend directory, this resolves correctly.
		exe, err := os.Executable()
		if err == nil {
			// Built binary: frontend is sibling of backend directory
			frontendDir = filepath.Join(filepath.Dir(exe), "..", "frontend")
		} else {
			frontendDir = "../frontend"
		}
	}

	mux := http.NewServeMux()

	// API endpoint
	mux.HandleFunc("/api/generate", corsMiddleware(handlers.GenerateCV))

	// Serve frontend static files
	mux.Handle("/", http.FileServer(http.Dir(frontendDir)))

	addr := ":" + port
	fmt.Printf("CV Generator server starting on http://localhost%s\n", addr)
	fmt.Printf("Open http://localhost%s in your browser to use the CV generator.\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// corsMiddleware adds CORS headers to allow requests from any origin
// (useful when the frontend is served separately during development).
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

