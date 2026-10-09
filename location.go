package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type LocationResult struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
}

type ElevationResult struct {
	Elevation []float64 `json:"elevation"`
}

func getCoordinates(location string) (float64, float64, error) {

	baseURL := "https://nominatim.openstreetmap.org/search"

	params := url.Values{}
	params.Set("q", location)
	params.Set("format", "json")
	params.Set("limit", "1")
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	request, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		return 0.0, 0.0, fmt.Errorf("request failed: %w", err)
	}

	request.Header.Set("User-Agent", contactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("response failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return 0.0, 0.0, fmt.Errorf("server returned HTTP status: %s", response.Status)
	}

	var results []LocationResult
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&results)

	if err != nil {
		return 0.0, 0.0, fmt.Errorf("failed to decode location data: %w", err)
	}

	if len(results) == 0 {
		return 0.0, 0.0, fmt.Errorf("no location data returned")
	}

	latitude, err := strconv.ParseFloat(results[0].Latitude, 64)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("failed to parse latitude: %w", err)
	}

	longitude, err := strconv.ParseFloat(results[0].Longitude, 64)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("failed to parse longitude: %w", err)
	}

	err = validateCoordinates(latitude, longitude)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("failed to validate coordinates: %w", err)
	}

	return latitude, longitude, nil
}

func parseCoordinates(coordinates string) (float64, float64, error) {
	parts := strings.Split(coordinates, ",")

	if len(parts) != 2 {
		return 0.0, 0.0, fmt.Errorf("coordinates must be in the form lat, long")
	}
	latitudeStr := strings.TrimSpace(parts[0])
	longitudeStr := strings.TrimSpace(parts[1])

	latitude, err := strconv.ParseFloat(latitudeStr, 64)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("invalid latitude %w", err)
	}
	longitude, err := strconv.ParseFloat(longitudeStr, 64)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("invalid longitude %w", err)
	}
	err = validateCoordinates(latitude, longitude)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("invalid coordinates: %w", err)
	}
	return latitude, longitude, nil
}

func validateCoordinates(latitude, longitude float64) error {
	if math.IsNaN(latitude) || math.IsInf(latitude, 0) {
		return fmt.Errorf("invalid latitude")
	}
	if latitude < -90 || latitude > 90 {
		return fmt.Errorf("latitude out of limits")
	}
	if math.IsNaN(longitude) || math.IsInf(longitude, 0) {
		return fmt.Errorf("invalid longitude")
	}
	if longitude < -180 || longitude > 180 {
		return fmt.Errorf("longitude out of limits")
	}
	return nil
}

func getElevation(latitude, longitude float64) (float64, error) {

	baseURL := "https://api.open-meteo.com/v1/elevation"

	params := url.Values{}
	params.Set("latitude", strconv.FormatFloat(latitude, 'f', 4, 64))
	params.Set("longitude", strconv.FormatFloat(longitude, 'f', 4, 64))
	requestURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	request, err := http.NewRequest("GET", requestURL, nil)

	if err != nil {
		return 0.0, fmt.Errorf("request failed: %w", err)
	}
	request.Header.Set("User-Agent", contactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0.0, fmt.Errorf("response failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return 0.0, fmt.Errorf("server returned HTTP status: %s", response.Status)
	}

	var result ElevationResult
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&result)

	if err != nil {
		return 0.0, fmt.Errorf("failed to decode elevation data: %w", err)
	}

	if len(result.Elevation) == 0 {
		return 0.0, fmt.Errorf("no elevation data returned")
	}

	return result.Elevation[0], nil
}
