package web

import (
	"groupie-tracker/data"
	"groupie-tracker/models"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// StartServer sets up routes, applies middleware, and starts the HTTP server
func StartServer() {
	// Fetch and cache all artist data
	allArtists = data.FetchFullArtistInfo()

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register static file handler
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Register route handlers
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/team", teamHandler)
	mux.HandleFunc("/artist/", artistHandler) // Dynamic handling inside function

	log.Println("🚀 Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", RecoveryMiddleware(mux)))
}

// Team page
func teamHandler(w http.ResponseWriter, _ *http.Request) {
	err := templates.ExecuteTemplate(w, "team.html", nil)
	if err != nil {
		log.Printf("🔥 Template error: %v", err)
		renderErrorPage(w, http.StatusInternalServerError, "team.html render error")
	}
}

// Render the homepage with pagination
func homeHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		renderErrorPage(w, http.StatusNotFound, "Page not found")
		return
	}

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

	err := templates.ExecuteTemplate(w, "index.html", viewData)
	if err != nil {
		log.Printf("🔥 Template error: %v", err)
		renderErrorPage(w, http.StatusInternalServerError, "index.html render error")
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
				log.Printf("🔥 Template error: %v", err)
				renderErrorPage(w, http.StatusInternalServerError, "artist.html render error")
			}
			return
		}
	}

	renderErrorPage(w, http.StatusNotFound, "Artist not found")
}
