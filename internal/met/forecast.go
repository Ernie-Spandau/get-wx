package met

import (
	"encoding/json"
	"fmt"
	"get-wx/internal/config"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type MetResponse struct {
	Properties Properties `json:"properties"`
}
type Properties struct {
	Meta       Meta         `json:"meta"`
	Timeseries []Timeseries `json:"timeseries"`
}
type Meta struct {
	Updated time.Time `json:"updated_at"`
}
type Timeseries struct {
	Data Data `json:"data"`
}
type Data struct {
	Instant     Instant  `json:"instant"`
	Next1Hours  Forecast `json:"next_1_hours"`
	Next6Hours  Forecast `json:"next_6_hours"`
	Next12Hours Forecast `json:"next_12_hours"`
}
type Instant struct {
	Details Details `json:"details"`
}
type Details struct {
	AirPressure    float64 `json:"air_pressure_at_sea_level"`
	AirTemperature float64 `json:"air_temperature"`
	CloudArea      float64 `json:"cloud_area_fraction"`
	Humidity       float64 `json:"relative_humidity"`
	WindDirection  float64 `json:"wind_from_direction"`
	WindSpeed      float64 `json:"wind_speed"`
}
type Forecast struct {
	Summary ForecastSummary `json:"summary"`
	Details ForecastDetails `json:"details"`
}
type ForecastSummary struct {
	SymbolCode string `json:"symbol_code"`
}
type ForecastDetails struct {
	PrecipitationAmount float64 `json:"precipitation_amount"`
}

func GetForecast(latitude, longitude, elevation float64) (MetResponse, error) {

	baseURL := "https://api.met.no/weatherapi/locationforecast/2.0/compact"

	params := url.Values{}
	params.Set("lat", strconv.FormatFloat(latitude, 'f', 4, 64))
	params.Set("lon", strconv.FormatFloat(longitude, 'f', 4, 64))
	params.Set("altitude", strconv.FormatFloat(elevation, 'f', 0, 64))
	url := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return MetResponse{}, fmt.Errorf("request failed %w", err)
	}

	request.Header.Set("User-Agent", config.ContactInfo)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return MetResponse{}, fmt.Errorf("response failed: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode > 299 {
		return MetResponse{}, fmt.Errorf("server returned HTTP status: %s", response.Status)
	}

	var forecast MetResponse
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&forecast)
	if err != nil {
		return MetResponse{}, fmt.Errorf("failed to decode weather data: %w", err)
	}
	if len(forecast.Properties.Timeseries) == 0 {
		return MetResponse{}, fmt.Errorf("no forecast data returned")
	}
	return forecast, nil
}
