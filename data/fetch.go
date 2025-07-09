package data

import (
	"encoding/json"
	"fmt"
	"groupie-tracker/models"
	"io"
	"log"
	"net/http"
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

func FetchFullArtistInfo() []models.FullArtistInfo {
	var artists []models.ArtistSummary
	if err := fetchJSON(baseURL+"/artists", &artists); err != nil {
		log.Fatal("Error fetching artists: ", err)
	}

	var fullData []models.FullArtistInfo

	for _, artist := range artists {
		var locations models.Locations
		var dates models.Dates
		var relation models.Relation

		if err := fetchJSON(artist.LocationsURL, &locations); err != nil {
			log.Println("Error fetching locations for artist ID", artist.ID, ":", err)
			continue
		}

		if err := fetchJSON(artist.DatesURL, &dates); err != nil {
			log.Println("Error fetching dates for artist ID", artist.ID, ":", err)
			continue
		}

		if err := fetchJSON(artist.RelationURL, &relation); err != nil {
			log.Println("Error fetching relation for artist ID", artist.ID, ":", err)
			continue
		}

		fullData = append(fullData, models.FullArtistInfo{
			ArtistSummary: artist,
			Locations:     locations.Locations,
			Dates:         dates.Dates,
			Relation:      relation.DatesLocations,
		})
	}

	return fullData
}
