package noaa

import (
	"encoding/json"
	"fmt"
	"get-wx/internal/config"
	"net/http"
)

type NOAAResponse struct {
	PointProperties Properties `json:"properties"`
}
type Properties struct {
	Stations string `json:"observationStations"`
}
type FeatureCollection struct {
	Features []Features `json:"features"`
}
type Features struct {
	ID string `json:"id"`
}

type ObservationResponse struct {
	Properties ObservationProperties `json:"properties"`
}
type ObservationProperties struct {
	Temperature               Measurement  `json:"temperature"`
	Dewpoint                  Measurement  `json:"dewpoint"`
	WindDirection             Measurement  `json:"windDirection"`
	WindSpeed                 Measurement  `json:"windSpeed"`
	WindGust                  Measurement  `json:"windGust"`
	BarometricPressure        Measurement  `json:"barometricPressure"`
	SeaLevelPressure          Measurement  `json:"seaLevelPressure"`
	Visibility                Measurement  `json:"visibility"`
	MaxTemperatureLast24Hours Measurement  `json:"maxTemperatureLast24Hours"`
	MinTemperatureLast24Hours Measurement  `json:"minTemperatureLast24Hours"`
	PrecipitationLast3Hours   Measurement  `json:"precipitationLast3Hours"`
	RelativeHumidity          Measurement  `json:"relativeHumidity"`
	WindChill                 Measurement  `json:"windChill"`
	HeatIndex                 Measurement  `json:"heatIndex"`
	CloudLayers               []CloudLayer `json:"cloudLayers"`
}
type Measurement struct {
	Value *float64 `json:"value"`
}
type CloudLayer struct {
	Base   Measurement `json:"base"`
	Amount string      `json:"amount"`
}

func GetStations(latitude, longitude float64) (FeatureCollection, error) {
	stationsURL, err := getStationsURL(latitude, longitude)
	if err != nil {
		fmt.Println("NOAA", err)
		return FeatureCollection{}, err
	}

	url := stationsURL.PointProperties.Stations

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return FeatureCollection{}, fmt.Errorf("request failed %w", err)
	}

	request.Header.Set("User-Agent", config.ContactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return FeatureCollection{}, fmt.Errorf("response failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return FeatureCollection{}, fmt.Errorf("server returned HTTP status: %s", response.Status)
	}

	var stations FeatureCollection
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&stations)
	if err != nil {
		return FeatureCollection{}, fmt.Errorf("failed to decode station data: %w", err)
	}

	// getStationWx(stations)
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

func GetStationWx(stations FeatureCollection) (ObservationResponse, error) {

	if len(stations.Features) == 0 {
		return ObservationResponse{}, fmt.Errorf("Failed to retrieve observation data from NOAA")
	}
	stationURL := stations.Features[0].ID + "/observations/latest"

	request, err := http.NewRequest("GET", stationURL, nil)
	if err != nil {
		return ObservationResponse{}, fmt.Errorf("request failed %w", err)
	}

	request.Header.Set("User-Agent", config.ContactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return ObservationResponse{}, fmt.Errorf("response failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return ObservationResponse{}, fmt.Errorf("server returned HTTP status: %s", response.Status)
	}

	var stationWx ObservationResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&stationWx)
	if err != nil {
		return ObservationResponse{}, fmt.Errorf("failed to decode station data: %w", err)
	}
	return stationWx, err
}
