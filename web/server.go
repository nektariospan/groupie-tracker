package web

import (
	"groupie-tracker/data"
	"groupie-tracker/models"
	"html/template"
	"log"
	"net/http"
	"runtime/debug"
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

// LoggingMiddleware logs each incoming request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

// RecoveryMiddleware recovers from panics and returns a 500 page
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v\n%s", err, debug.Stack())
				renderErrorPage(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// StartServer sets up routes, applies middleware, and starts the HTTP server
func StartServer() {
	// Fetch and cache all artist data
	allArtists = data.FetchFullArtistInfo()

	// Register static file handler and central router
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", http.HandlerFunc(router))

	// Wrap the default mux with logging and recovery middleware
	var handler http.Handler = http.DefaultServeMux
	handler = LoggingMiddleware(handler)
	handler = RecoveryMiddleware(handler)

	log.Println("🚀 Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// Central router for handling known and unknown routes
func router(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		homeHandler(w, r)
	case strings.HasPrefix(r.URL.Path, "/artist/"):
		artistHandler(w, r)
	case r.URL.Path == "/team":
		teamHandler(w, r)
	default:
		renderErrorPage(w, http.StatusNotFound, "Page not found")
	}
}

func teamHandler(w http.ResponseWriter, _ *http.Request) {
	if err := templates.ExecuteTemplate(w, "team.html", nil); err != nil {
		renderErrorPage(w, http.StatusInternalServerError, "Team page template error")
	}
}

// Render the homepage with pagination
func homeHandler(w http.ResponseWriter, r *http.Request) {
	total := len(allArtists)

	// Default values
	limit := 10
	page := 1

	// Handle ?limit=all or numeric values
	limitParam := r.URL.Query().Get("limit")
	switch limitParam {
	case "all":
		limit = total
	case "":
		limitParam = strconv.Itoa(limit)
	default:
		if parsed, err := strconv.Atoi(limitParam); err == nil && parsed > 0 {
			limit = parsed
		} else {
			limitParam = strconv.Itoa(limit)
		}
	}

	// Handle ?page
	if pageParam := r.URL.Query().Get("page"); pageParam != "" {
		if parsed, err := strconv.Atoi(pageParam); err == nil && parsed > 0 {
			page = parsed
		}
	}

	// Compute slice boundaries
	start := (page - 1) * limit
	end := start + limit
	if start >= total {
		start, end, page = 0, limit, 1
	}
	if end > total {
		end = total
	}

	// Prepare view data
	viewData := struct {
		Artists      []models.FullArtistInfo
		Limit        string
		Page         int
		Total        int
		HasNext      bool
		HasPrevious  bool
		NextPage     int
		PreviousPage int
	}{
		Artists:      allArtists[start:end],
		Limit:        limitParam,
		Page:         page,
		Total:        total,
		HasNext:      end < total,
		HasPrevious:  start > 0,
		NextPage:     page + 1,
		PreviousPage: page - 1,
	}

	if err := templates.ExecuteTemplate(w, "index.html", viewData); err != nil {
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
			if err := templates.ExecuteTemplate(w, "artist.html", artist); err != nil {
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
	if err := errorTemplate.Execute(w, ErrorData{Code: code, Message: message}); err != nil {
		http.Error(w, "Error rendering error page", http.StatusInternalServerError)
	}
}
