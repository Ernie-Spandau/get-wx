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

const contactInfo = "get-wx/0.1 https://github.com/Ernie-Spandau/get-wx"

type LocationResult struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"lon"`
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

type Data struct {
	Instant Instant `json:"instant"`
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
			fmt.Println("Failed to find location")
			return
		}

		fmt.Println("Weather for", *location)

	} else {

		parts := strings.Split(*coordinates, ",")
		if len(parts) != 2 {
			fmt.Println("Coordinates must be in the form lat, long")
			return
		}

		latitude, err = strconv.ParseFloat(parts[0], 64)
		if err != nil || math.IsNaN(latitude) || math.IsInf(latitude, 0) {
			fmt.Println("Invalid latitude")
			return
		}
		if latitude < -90 || latitude > 90 {
			fmt.Println("Latitude out of limits (-90° - 90°)")
			return
		}

		longitude, err = strconv.ParseFloat(parts[1], 64)
		if err != nil || math.IsNaN(longitude) || math.IsInf(longitude, 0) {
			fmt.Println("Invalid longitude")
			return
		}
		if longitude < -180 || longitude > 180 {
			fmt.Println("Longitude out of limits (-180° - 180°)")
			return
		}
		fmt.Printf("Weather for %.4f, %.4f\n", latitude, longitude)
	}
	getWeather(latitude, longitude)
}

func getWeather(latitude, longitude float64) {
	url := fmt.Sprintf("https://api.met.no/weatherapi/locationforecast/2.0/compact?lat=%.4f&lon=%.4f", latitude, longitude)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Request failed")
		return
	}

	request.Header.Set("User-Agent", contactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Response failed")
		return
	}

	defer response.Body.Close()
	var weather WeatherResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&weather)
	if err != nil {
		fmt.Println("Failed to decode weather data")
		return
	}
	if len(weather.Properties.Timeseries) == 0 {
		fmt.Println("No forecast data returned")
		return
	}

	wxData := weather.Properties.Timeseries[0].Data.Instant.Details
	temperature := wxData.AirTemperature
	temperatureF := (temperature * 9.0 / 5.0) + 32
	pressure := wxData.AirPressure
	pressureInHg := pressure * 0.02953
	windDirection := int(math.Round(wxData.WindDirection))
	windspeed := wxData.WindSpeed
	humidity := int(math.Round(wxData.Humidity))
	skyCondition := wxData.CloudArea
	fmt.Printf("Current temperature is: %.1f°C, %.1f°F\n", temperature, temperatureF)
	fmt.Printf("Current air pressure is %.2fhPa, %.2finHg\n", pressure, pressureInHg)
	fmt.Printf("Current wind direction is: %v°\n", windDirection)
	fmt.Println("Current wind speed is:", windspeed)
	fmt.Printf("Current humidity is: %d%%\n", humidity)
	fmt.Printf("Current sky condition is %.0f%% cloudy", skyCondition)
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
		fmt.Println("Request failed")
		return 0.0, 0.0, err
	}

	request.Header.Set("User-Agent", contactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Response failed")
		return 0.0, 0.0, err
	}

	defer response.Body.Close()

	var results []LocationResult
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&results)

	if err != nil {
		fmt.Println("Failed to decode location data")
		return 0.0, 0.0, err
	}
	if len(results) == 0 {
		fmt.Println("No location data returned")
		return 0.0, 0.0, fmt.Errorf("No location data returned")
	}

	latitude, err := strconv.ParseFloat(results[0].Latitude, 64)
	if err != nil || math.IsNaN(latitude) || math.IsInf(latitude, 0) {
		fmt.Println("Invalid latitude")
		return 0.0, 0.0, err
	}
	if latitude < -90 || latitude > 90 {
		fmt.Println("Latitude out of limits (-90° - 90°)")
		return 0.0, 0.0, fmt.Errorf("Latitude out of limits")
	}

	longitude, err := strconv.ParseFloat(results[0].Longitude, 64)
	if err != nil || math.IsNaN(longitude) || math.IsInf(longitude, 0) {
		fmt.Println("Invalid longitude")
		return 0.0, 0.0, err
	}
	if longitude < -180 || longitude > 180 {
		fmt.Println("Longitude out of limits (-180° - 180°)")
		return 0.0, 0.0, fmt.Errorf("Longitude out of limits")
	}

	return latitude, longitude, nil
}
