package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

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
		ForecastCondition1H:      formatCondition(next1HrData.Summary.SymbolCode),
		ForecastCondition6H:      formatCondition(next6HrData.Summary.SymbolCode),
		ForecastCondition12H:     formatCondition(next12HrData.Summary.SymbolCode),
		ForecastPrecipitation1H:  next1HrData.Details.PrecipitationAmount,
		ForecastPrecipitation6H:  next6HrData.Details.PrecipitationAmount,
		ForecastPrecipitation12H: next12HrData.Details.PrecipitationAmount,
	}
	return wx
}
