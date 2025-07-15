package data

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"io"
	"log"
	"net/http"
	"sync"
)

const baseURL = "https://groupietrackers.herokuapp.com/api"

// fetchJSON performs an HTTP GET to the given URL and decodes the JSON response into target.
func fetchJSON(url string, target any) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET error: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Read error: %v", err)
	}

	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("JSON decode error: %v", err)
	}
	return nil
}

// FetchFullArtistInfo fetches summary list and then, in parallel for each artist,
// retrieves locations, dates and relations before returning the combined slice.
func FetchFullArtistInfo() []models.FullArtistInfo {
	// First, fetch the list of all artist summaries.
	var artists []models.ArtistSummary
	if err := fetchJSON(baseURL+"/artists", &artists); err != nil {
		log.Fatal("Error fetching artists: ", err)
	}

	// Prepare synchronization primitives
	var wg sync.WaitGroup
	artistChan := make(chan models.FullArtistInfo, len(artists))

	// Launch one goroutine per artist to fetch details concurrently
	for _, artist := range artists {
		wg.Add(1)

		go func(artist models.ArtistSummary) {
			defer wg.Done()

			// Containers for each sub-request
			var (
				locations models.Locations
				dates     models.Dates
				relation  models.Relation
			)

			// innerWg waits for the three fetches
			var innerWg sync.WaitGroup
			errChan := make(chan error, 3) // buffered to hold up to 3 errors

			// Fetch locations concurrently
			innerWg.Add(1)
			go func() {
				defer innerWg.Done()
				if err := fetchJSON(artist.LocationsURL, &locations); err != nil {
					errChan <- fmt.Errorf("locations: %v", err)
				}
			}()

			// Fetch dates concurrently
			innerWg.Add(1)
			go func() {
				defer innerWg.Done()
				if err := fetchJSON(artist.DatesURL, &dates); err != nil {
					errChan <- fmt.Errorf("dates: %v", err)
				}
			}()

			// Fetch relation concurrently
			innerWg.Add(1)
			go func() {
				defer innerWg.Done()
				if err := fetchJSON(artist.RelationURL, &relation); err != nil {
					errChan <- fmt.Errorf("relation: %v", err)
				}
			}()

			// Wait for all three to complete
			innerWg.Wait()
			close(errChan)

			// If any fetch failed, log and skip this artist
			for e := range errChan {
				log.Println("Error fetching for artist ID", artist.ID, ":", e)
				return
			}

			// All good: send the aggregated FullArtistInfo into the channel
			artistChan <- models.FullArtistInfo{
				ArtistSummary: artist,
				Locations:     locations.Locations,
				Dates:         dates.Dates,
				Relation:      relation.DatesLocations,
			}

		}(artist)
	}

	// Close the channel once all artist goroutines have finished
	go func() {
		wg.Wait()
		close(artistChan)
	}()

	// Collect results into a slice
	var fullData []models.FullArtistInfo
	for a := range artistChan {
		fullData = append(fullData, a)
	}

	return fullData
}
