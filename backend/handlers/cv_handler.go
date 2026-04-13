package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/screenleon/cv-generate/generators"
	"github.com/screenleon/cv-generate/models"
)

// GenerateCV handles POST /api/generate requests.
// It reads JSON CV data from the request body, generates the document,
// and responds with the document as a downloadable file.
func GenerateCV(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data models.CVData
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "http: request body too large" {
			status = http.StatusRequestEntityTooLarge
		}
		http.Error(w, "Invalid request body: "+err.Error(), status)
		return
	}

	// Validate required fields
	if strings.TrimSpace(data.Name) == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Default values
	if data.Template == "" {
		data.Template = "simple"
	}
	if data.Format == "" {
		data.Format = "pdf"
	}

	var (
		docBytes    []byte
		err         error
		contentType string
		filename    string
	)

	switch strings.ToLower(data.Format) {
	case "word", "docx":
		docBytes, err = generators.GenerateWord(&data)
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		filename = "cv.docx"
	default: // "pdf"
		docBytes, err = generators.GeneratePDF(&data)
		contentType = "application/pdf"
		filename = "cv.pdf"
	}

	if err != nil {
		http.Error(w, "Failed to generate document: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(docBytes)
}
