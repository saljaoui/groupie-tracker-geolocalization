package Groupie_tracker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var location Location

// func to Get All data From json
func GetArtistsDataStruct() ([]JsonData, error) {
	var artistData []JsonData

	response, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		return nil, fmt.Errorf("error fetching data from artist data: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching data from artist data: %d", response.StatusCode)
	}

	err = json.NewDecoder(response.Body).Decode(&artistData)
	if err != nil {
		return nil, fmt.Errorf("error fetching data from artist data: %v", err)
	}

	return artistData, nil
}

// / func to fetching data from any struct and return Struct Artist with Id user
func FetchDataRelationFromId(id string) (Artist, error) {
	url := "https://groupietrackers.herokuapp.com/api"
	var artist Artist
	err := GetanyStruct(url+"/artists/"+id, &artist)
	if err != nil {
		return Artist{}, fmt.Errorf("error fetching data from artist data: %w", err)
	}

	if artist.ID == 0 {
		return artist, nil
	}

	var date Date
	err = GetanyStruct(url+"/dates/"+id, &date)
	if err != nil {
		return Artist{}, fmt.Errorf("error fetching data from artist data: %w", err)
	}

	err = GetanyStruct(url+"/locations/"+id, &location)
	if err != nil {
		return Artist{}, fmt.Errorf("error fetching data from locations data: %w", err)
	}

	var relation Relation
	err = GetanyStruct(url+"/relation/"+id, &relation)
	if err != nil {
		return Artist{}, fmt.Errorf("error fetching data from locations data: %w", err)
	}
	artist.Location = location.Location
	artist.Date = date.Date

	artist.DatesLocations = formatLocations(relation.DatesLocations)

	return artist, nil
}

type GeocodeResponse struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

func GeoMap(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	indexStr := r.URL.Query().Get("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(location.Location) {
		http.Error(w, "Invalid index", http.StatusBadRequest)
		return
	}
	result, err := geocodeAddress(location.Location[index])
	if err != nil {
		http.Error(w, "Error geocoding address", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}

func geocodeAddress(address string) (string, error) {
	var result string
	encodedAddress := url.QueryEscape(address)
	url := fmt.Sprintf("https://nominatim.openstreetmap.org/search?format=json&q=%s", encodedAddress)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error: %s", resp.Status)
	}

	var results []GeocodeResponse
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return "", fmt.Errorf("no results the error is: %s", err)
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no results found for address: %s", address)
	}
	result = fmt.Sprintf("https://www.google.com/maps/place/%v,%v/@%v,%v,7778m/", results[0].Lat, results[0].Lon, results[0].Lat, results[0].Lon)

	return result, nil
}

func GetanyStruct(url string, result interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error fetching data from URL: %w", err)
	}
	defer response.Body.Close()
	// Decode the JSON response into the provided struct
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return fmt.Errorf("error decoding JSON data: %w", err)
	}
	return nil
}

// func To Format String To remove '_' or '-' and Capitaliz text
func formatLocations(locations map[string][]string) map[string][]string {
	formatted := make(map[string][]string, len(locations))

	for location, dates := range locations {
		formattedLoc := (strings.NewReplacer("-", " ", "_", " ").Replace(location))
		formatted[formattedLoc] = dates
	}

	return formatted
}
