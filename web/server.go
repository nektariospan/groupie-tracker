package web

import (
	"groupie-tracker/data"
	"groupie-tracker/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Templates
var templates = template.Must(template.ParseGlob("templates/*.html"))
var errorTemplate = template.Must(template.ParseFiles("templates/error.html"))

// Cached data
var allArtists []models.FullArtistInfo

// Error data for rendering error.html
type ErrorData struct {
	Code    int
	Message string
}

// Start the server
func StartServer() {
	allArtists = data.FetchFullArtistInfo()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static")))) // for static files
	http.Handle("/", http.HandlerFunc(router))                                                 // central router

	log.Println("🚀 Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Central router for handling known and unknown routes
func router(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		homeHandler(w, r)
	case strings.HasPrefix(r.URL.Path, "/artist/"):
		artistHandler(w, r)
	default:
		renderErrorPage(w, http.StatusNotFound, "Page not found")
	}
}

// Render the homepage
func homeHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", allArtists)
	if err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Home page template error")
	}
}

// Render individual artist page
func artistHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/artist/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		renderErrorPage(w, http.StatusBadRequest, "Invalid artist ID in URL")
		return
	}

	for _, artist := range allArtists {
		if artist.ID == id {
			err := templates.ExecuteTemplate(w, "artist.html", artist)
			if err != nil {
				renderErrorPage(w, http.StatusInternalServerError, "Error rendering artist page")
			}
			return
		}
	}

	renderErrorPage(w, http.StatusNotFound, "Artist not found")
}

// Renders error.html with given code and message
func renderErrorPage(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	err := errorTemplate.Execute(w, ErrorData{Code: code, Message: message})
	if err != nil {
		http.Error(w, "Error rendering error page", http.StatusInternalServerError)
	}
}
