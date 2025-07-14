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

func fetchJSON(url string, target interface{}) error {
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

// ΑΣΥΓΧΡΟΝΗ ΣΥΝΑΡΤΗΣΗ
func FetchFullArtistInfo() []models.FullArtistInfo {
	var artists []models.ArtistSummary
	if err := fetchJSON(baseURL+"/artists", &artists); err != nil {
		log.Fatal("Error fetching artists: ", err)
	}

	var wg sync.WaitGroup
	artistChan := make(chan models.FullArtistInfo, len(artists))

	for _, artist := range artists {
		wg.Add(1)

		go func(artist models.ArtistSummary) {
			defer wg.Done()

			var (
				locations models.Locations
				dates     models.Dates
				relation  models.Relation
			)

			var innerWg sync.WaitGroup
			errChan := make(chan error, 3)

			// Fetch Locations
			innerWg.Add(1)
			go func() {
				defer innerWg.Done()
				if err := fetchJSON(artist.LocationsURL, &locations); err != nil {
					errChan <- fmt.Errorf("locations: %v", err)
				}
			}()

			// Fetch Dates
			innerWg.Add(1)
			go func() {
				defer innerWg.Done()
				if err := fetchJSON(artist.DatesURL, &dates); err != nil {
					errChan <- fmt.Errorf("dates: %v", err)
				}
			}()

			// Fetch Relation
			innerWg.Add(1)
			go func() {
				defer innerWg.Done()
				if err := fetchJSON(artist.RelationURL, &relation); err != nil {
					errChan <- fmt.Errorf("relation: %v", err)
				}
			}()

			// Wait for all 3 to finish
			innerWg.Wait()
			close(errChan)

			// Check if there were any errors
			for e := range errChan {
				log.Println("Error fetching for artist ID", artist.ID, ":", e)
				return
			}

			artistChan <- models.FullArtistInfo{
				ArtistSummary: artist,
				Locations:     locations.Locations,
				Dates:         dates.Dates,
				Relation:      relation.DatesLocations,
			}
		}(artist)
	}

	// Κλείνουμε το κανάλι όταν τελειώσουν όλα
	go func() {
		wg.Wait()
		close(artistChan)
	}()

	var fullData []models.FullArtistInfo
	for a := range artistChan {
		fullData = append(fullData, a)
	}

	return fullData
}
