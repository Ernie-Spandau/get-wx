package main

import (
	"encoding/json"
	"flag"
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

type Details struct {
	AirPressure    float64 `json:"air_pressure_at_sea_level"`
	AirTemperature float64 `json:"air_temperature"`
	CloudArea      float64 `json:"cloud_area_fraction"`
	Humidity       float64 `json:"relative_humidity"`
	WindDirection  float64 `json:"wind_from_direction"`
	WindSpeed      float64 `json:"wind_speed"`
}

type Instant struct {
	Details Details `json:"details"`
}

type ForecastSummary struct {
	SymbolCode string `json:"symbol_code"`
}

type ForecastDetails struct {
	PrecipitationAmount float64 `json:"precipitation_amount"`
}

type Forecast struct {
	Summary ForecastSummary `json:"summary"`
	Details ForecastDetails `json:"details"`
}

type Data struct {
	Instant     Instant  `json:"instant"`
	Next1Hours  Forecast `json:"next_1_hours"`
	Next6Hours  Forecast `json:"next_6_hours"`
	Next12Hours Forecast `json:"next_12_hours"`
}

type Timeseries struct {
	Data Data `json:"data"`
}

type Properties struct {
	Timeseries []Timeseries `json:"timeseries"`
}

type WeatherResponse struct {
	Properties Properties `json:"properties"`
}

type Weather struct {
	Temperature              float64
	TemperatureF             float64
	Pressure                 float64
	PressureInHg             float64
	WindDirection            int
	WindSpeed                float64
	WindSpeedMPH             float64
	WindSpeedKTS             float64
	Humidity                 int
	SkyCondition             int
	ForecastCondition1H      string
	ForecastPrecipitation1H  float64
	ForecastCondition6H      string
	ForecastPrecipitation6H  float64
	ForecastCondition12H     string
	ForecastPrecipitation12H float64
	Elevation                float64
	ElevationFT              float64
}

const contactInfo = "get-wx/0.1 https://github.com/Ernie-Spandau/get-wx"

func main() {

	location := flag.String("l", "", "location for weather")
	coordinates := flag.String("c", "", "coordinates as lattitude, longitude")
	var latitude, longitude float64
	var err error

	flag.Parse()
	if *location == "" && *coordinates == "" {
		fmt.Println("Missing location or coordinates")
		return
	}
	if *location != "" && *coordinates != "" {
		fmt.Println("Use either -l or -c, not both")
		return
	}
	if *location != "" {
		latitude, longitude, err = getCoordinates(*location)
		if err != nil {
			fmt.Printf("Failed to find location: %v\n", err)
			return
		}

		fmt.Printf("Weather for %s (%.4f, %.4f)\n", *location, latitude, longitude)

	} else {

		latitude, longitude, err = parseCoordinates(*coordinates)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Weather for %.4f, %.4f\n", latitude, longitude)
	}
	elevation, err := getElevation(latitude, longitude)
	if err != nil {
		fmt.Println(err)
		return
	}

	weather, err := getWeather(latitude, longitude, elevation)
	if err != nil {
		fmt.Println(err)
		return
	}

	wx := convertWeather(weather, elevation)
	displayWeather(wx)
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

func getWeather(latitude, longitude, elevation float64) (WeatherResponse, error) {

	baseURL := "https://api.met.no/weatherapi/locationforecast/2.0/compact"

	params := url.Values{}
	params.Set("lat", strconv.FormatFloat(latitude, 'f', 4, 64))
	params.Set("lon", strconv.FormatFloat(longitude, 'f', 4, 64))
	params.Set("altitude", strconv.FormatFloat(elevation, 'f', 0, 64))
	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("request failed %w", err)
	}

	request.Header.Set("User-Agent", contactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("response failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return WeatherResponse{}, fmt.Errorf("server returned HTTP status: %s", response.Status)
	}

	var weather WeatherResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&weather)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("failed to decode weather data: %w", err)
	}
	if len(weather.Properties.Timeseries) == 0 {
		return WeatherResponse{}, fmt.Errorf("no forecast data returned")
	}
	return weather, nil
}

func convertWeather(weatherResponse WeatherResponse, elevation float64) Weather {

	wxData := weatherResponse.Properties.Timeseries[0].Data.Instant.Details
	next1HrData := weatherResponse.Properties.Timeseries[0].Data.Next1Hours
	next6HrData := weatherResponse.Properties.Timeseries[0].Data.Next6Hours
	next12HrData := weatherResponse.Properties.Timeseries[0].Data.Next12Hours

	wx := Weather{
		Elevation:                elevation,
		ElevationFT:              elevation * 3.28084,
		Temperature:              wxData.AirTemperature,
		TemperatureF:             (wxData.AirTemperature * 9.0 / 5.0) + 32,
		Pressure:                 wxData.AirPressure,
		PressureInHg:             wxData.AirPressure * 0.02953,
		WindDirection:            int(math.Round(wxData.WindDirection)),
		WindSpeed:                wxData.WindSpeed,
		WindSpeedMPH:             wxData.WindSpeed * 2.2369362921,
		WindSpeedKTS:             wxData.WindSpeed * 1.9438444924,
		Humidity:                 int(math.Round(wxData.Humidity)),
		SkyCondition:             int(math.Round(wxData.CloudArea)),
		ForecastCondition1H:      next1HrData.Summary.SymbolCode,
		ForecastCondition6H:      next6HrData.Summary.SymbolCode,
		ForecastCondition12H:     next12HrData.Summary.SymbolCode,
		ForecastPrecipitation1H:  next1HrData.Details.PrecipitationAmount,
		ForecastPrecipitation6H:  next6HrData.Details.PrecipitationAmount,
		ForecastPrecipitation12H: next12HrData.Details.PrecipitationAmount,
	}
	return wx
}

func displayWeather(wx Weather) {
	fmt.Printf("Elevation: %.2f m %.2f ft MSL\n", wx.Elevation, wx.ElevationFT)
	fmt.Printf("Current temperature is: %.1f°C, %.1f°F\n", wx.Temperature, wx.TemperatureF)
	fmt.Printf("Current air pressure is %.2fhPa, %.2finHg\n", wx.Pressure, wx.PressureInHg)
	fmt.Printf("Current wind speed is: %.1f m/s, %.1f MpH, %.1f KTS\n", wx.WindSpeed, wx.WindSpeedMPH, wx.WindSpeedKTS)
	fmt.Printf("Current wind direction is: %03d°\n", wx.WindDirection)
	fmt.Printf("Current humidity is: %d%%\n", wx.Humidity)
	fmt.Printf("Current sky condition is: %d%% cloudy\n\n", wx.SkyCondition)

	fmt.Println("Forecast for the next hour:")
	fmt.Printf("Sky condition: %v\n", wx.ForecastCondition1H)
	fmt.Printf("Preciptitation %.1f in\n\n", wx.ForecastPrecipitation1H/25.4)

	fmt.Println("Forecast for the next 6 hours:")
	fmt.Printf("Sky condition: %v\n", wx.ForecastCondition6H)
	fmt.Printf("Preciptitation: %.1f in\n\n", wx.ForecastPrecipitation6H/25.4)

	fmt.Println("Forecast for the next 12 hours:")
	fmt.Printf("Sky condition: %v\n", wx.ForecastCondition12H)
	fmt.Printf("Preciptitation: %.1f in\n\n", wx.ForecastPrecipitation12H/25.4)

}
