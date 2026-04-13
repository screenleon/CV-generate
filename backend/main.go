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

	// Resolve frontend directory: allow override via env, otherwise prefer
	// working-directory-relative locations before falling back to an
	// executable-relative path for built binaries.
	frontendDir := os.Getenv("FRONTEND_DIR")
	if frontendDir == "" {
		candidates := []string{
			"../frontend",
			"./frontend",
		}
		for _, candidate := range candidates {
			if info, err := os.Stat(candidate); err == nil && info.IsDir() {
				frontendDir = candidate
				break
			}
		}
		if frontendDir == "" {
			exe, err := os.Executable()
			if err == nil {
				candidate := filepath.Join(filepath.Dir(exe), "..", "frontend")
				if info, err := os.Stat(candidate); err == nil && info.IsDir() {
					frontendDir = candidate
				}
			}
		}
		if frontendDir == "" {
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

// corsMiddleware adds CORS headers. The allowed origin defaults to same-origin
// (empty, browsers will block cross-origin requests) unless CORS_ORIGIN is set.
// Set CORS_ORIGIN=* during development to allow all origins.
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	origin := os.Getenv("CORS_ORIGIN")
	return func(w http.ResponseWriter, r *http.Request) {
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

