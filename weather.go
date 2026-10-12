package main

import (
	"get-wx/internal/met"
	"math"
	"time"
)

type Weather struct {
	Updated                  time.Time
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

func convertWeather(metResponse met.MetResponse, elevation float64) Weather {
	wxUpdatedTime := metResponse.Properties.Meta
	wxData := metResponse.Properties.Timeseries[0].Data.Instant.Details
	next1HrData := metResponse.Properties.Timeseries[0].Data.Next1Hours
	next6HrData := metResponse.Properties.Timeseries[0].Data.Next6Hours
	next12HrData := metResponse.Properties.Timeseries[0].Data.Next12Hours

	wx := Weather{
		Updated:                  wxUpdatedTime.Updated,
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
