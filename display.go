package main

import (
	"fmt"
	"get-wx/internal/met"
	"get-wx/internal/noaa"
	"time"
)

func displayMeasurement(label string, value *float64, format string) {
	if value != nil {
		fmt.Printf("%-32s ", label)
		fmt.Printf(format+"\n", *value)
	} else {
		fmt.Printf("%-32s No data\n", label)
	}
}

func displayDate(label string, date time.Time) {
	fmt.Printf("%-32s %s\n",
		label,
		date.Local().Format("2006-01-02 15:04 MST"),
	)
}

func displayForecast(wx met.Weather) {
	fmt.Printf("\n----------------------------------\n")
	displayDate("Last Updated:", wx.Updated)
	displayMeasurement("Elevation:", &wx.ElevationFT, "%.2f ft")
	displayMeasurement("Temperature:", &wx.Temperature, "%.1f° C")
	displayMeasurement("Temperature:", &wx.TemperatureF, "%.1f° F")
	displayMeasurement("Baromatric pressure:", &wx.PressureInHg, "%.2f inHg")
	displayMeasurement("Wind speed:", &wx.WindSpeed, "%.1f m/s")
	displayMeasurement("Wind speed:", &wx.WindSpeedMPH, "%.1f MpH")
	displayMeasurement("Wind speed:", &wx.WindSpeedKTS, "%.1f KTS")
	displayMeasurement("Wind direction:", &wx.WindDirection, "%3.0f°")
	displayMeasurement("humidity:", &wx.Humidity, "%.1f%%")
	displayMeasurement("Sky condition:", &wx.SkyCondition, "%.1f%%")
	fmt.Printf("\n")
	fmt.Println("Forecast for the next hour:")
	fmt.Printf("%-32s %s\n", "Sky condition:", wx.ForecastCondition1H)
	displayMeasurement("Preciptitation:", &wx.ForecastPrecipitation1H, "%.1f")
	fmt.Printf("\n")
	fmt.Println("Forecast for the next 6 hours:")
	fmt.Printf("%-32s %s\n", "Sky condition:", wx.ForecastCondition6H)
	displayMeasurement("Preciptitation:", &wx.ForecastPrecipitation6H, "%.1f")
	fmt.Printf("\n")
	fmt.Println("Forecast for the next 12 hours:")
	fmt.Printf("%-32s %s\n", "Sky condition:", wx.ForecastCondition12H)
}

func displayObservations(observations noaa.ObservationResponse) {
	obs := observations.Properties
	fmt.Printf("\n----------------------------------\n")
	fmt.Println("Current NOAA measured observations at", obs.StationID, obs.StationName)
	displayDate("Last observation:", obs.Timestamp)
	displayMeasurement("Temperature:", obs.Temperature.Value, "%.1f° C")
	displayMeasurement("Dewpoint:", obs.Dewpoint.Value, "%.1f° C")
	displayMeasurement("Wind Direction:", obs.WindDirection.Value, "%3.0f°")
	displayMeasurement("Wind speed:", obs.WindSpeed.Value, "%.1f KpH")
	displayMeasurement("Wind gust:", obs.WindGust.Value, "%.1f KpH")
	displayMeasurement("Barometric Pressure:", obs.BarometricPressure.Value, "%.2f Pa")
	displayMeasurement("Sea Level Pressure:", obs.SeaLevelPressure.Value, "%.2f Pa")
	displayMeasurement("Visibility:", obs.Visibility.Value, "%.1f m")
	displayMeasurement("Max Temperature Last 24 Hours:", obs.MaxTemperatureLast24Hours.Value, "%.1f° C")
	displayMeasurement("Min Temperature Last 24 Hours:", obs.MinTemperatureLast24Hours.Value, "%.1f° C")
	displayMeasurement("Precipitation Last 3 Hours:", obs.PrecipitationLast3Hours.Value, "%.1f mm")
	displayMeasurement("Relative Humidity:", obs.RelativeHumidity.Value, "%.1f%%")
	displayMeasurement("Windchill:", obs.WindChill.Value, "%.2f° C")
	displayMeasurement("Heat Index:", obs.HeatIndex.Value, "%.2f° C")
}
