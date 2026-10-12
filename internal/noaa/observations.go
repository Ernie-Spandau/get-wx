package noaa

import (
	"encoding/json"
	"fmt"
	"get-wx/internal/config"
	"net/http"
)

type NOAAResponse struct {
	Properties Properties `json:"properties"`
}
type Properties struct {
	Stations string `json:"observationStations"`
}
type FeatureCollection struct {
	Features Features `json:"features"`
}
type Features struct {
	ID []string `json:"id"`
}

func GetStations(latitude, longitude float64) (Features, error) {
	stationsURL, err := getStationsURL(latitude, longitude)
	if err != nil {
		fmt.Println("NOAA", err)
		return Features{}, err
	}

	url := stationsURL.Properties.Stations

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Features{}, fmt.Errorf("request failed %w", err)
	}

	request.Header.Set("User-Agent", config.ContactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return Features{}, fmt.Errorf("response failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return Features{}, fmt.Errorf("server returned HTTP status: %s", response.Status)
	}

	var stations Features
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&stationsURL)
	if err != nil {
		return Features{}, fmt.Errorf("failed to decode station data: %w", err)
	}
	return stations, nil
}

func getStationsURL(latitude, longitude float64) (NOAAResponse, error) {

	baseURL := "https://api.weather.gov/points/"

	url := fmt.Sprintf("%s%.4f,%.4f", baseURL, latitude, longitude)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return NOAAResponse{}, fmt.Errorf("request failed %w", err)
	}

	request.Header.Set("User-Agent", config.ContactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return NOAAResponse{}, fmt.Errorf("response failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return NOAAResponse{}, fmt.Errorf("server returned HTTP status: %s", response.Status)
	}

	var stationsURL NOAAResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&stationsURL)
	if err != nil {
		return NOAAResponse{}, fmt.Errorf("failed to decode weather data: %w", err)
	}

	return stationsURL, nil
}
