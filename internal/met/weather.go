package met

import (
	"time"
)

type Weather struct {
	Updated                  time.Time
	Temperature              float64
	TemperatureF             float64
	Pressure                 float64
	PressureInHg             float64
	WindDirection            float64
	WindSpeed                float64
	WindSpeedMPH             float64
	WindSpeedKTS             float64
	Humidity                 float64
	SkyCondition             float64
	ForecastCondition1H      string
	ForecastPrecipitation1H  float64
	ForecastCondition6H      string
	ForecastPrecipitation6H  float64
	ForecastCondition12H     string
	ForecastPrecipitation12H float64
	Elevation                float64
	ElevationFT              float64
}

func ConvertWeather(metResponse MetResponse, elevation float64) Weather {
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
		WindDirection:            wxData.WindDirection,
		WindSpeed:                wxData.WindSpeed,
		WindSpeedMPH:             wxData.WindSpeed * 2.2369362921,
		WindSpeedKTS:             wxData.WindSpeed * 1.9438444924,
		Humidity:                 wxData.Humidity,
		SkyCondition:             wxData.CloudArea,
		ForecastCondition1H:      formatCondition(next1HrData.Summary.SymbolCode),
		ForecastCondition6H:      formatCondition(next6HrData.Summary.SymbolCode),
		ForecastCondition12H:     formatCondition(next12HrData.Summary.SymbolCode),
		ForecastPrecipitation1H:  next1HrData.Details.PrecipitationAmount / 25.4,
		ForecastPrecipitation6H:  next6HrData.Details.PrecipitationAmount / 25.4,
		ForecastPrecipitation12H: next12HrData.Details.PrecipitationAmount / 25.4,
	}
	return wx
}
