package web

import (
	"bytes"
	"groupie-tracker/models"
	"html/template"
	"log"
	"net/http"
)

// Templates
var templates = template.Must(template.ParseGlob("templates/*.html"))
var errorTemplate = template.Must(template.ParseFiles("templates/errors/error.html"))

// Cached data
var allArtists []models.FullArtistInfo

// Error data for rendering error.html
type ErrorData struct {
	Code    int
	Message string
}

// renderErrorPage writes the error template using a buffer to avoid
// sending headers prematurely. Only one WriteHeader is called.
func renderErrorPage(w http.ResponseWriter, code int, message string) {
	var buf bytes.Buffer
	data := ErrorData{Code: code, Message: message}

	// execute template into buffer
	if err := errorTemplate.Execute(&buf, data); err != nil {
		log.Printf("Error rendering error page template: %v", err)
		// fallback error if template execution fails
		http.Error(w, "An unexpected error occurred.", http.StatusInternalServerError)
		return
	}

	// write headers and content
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("Client disconnected (broken pipe): %v", err)
	}
}
