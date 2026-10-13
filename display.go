package main

import (
	"fmt"
	"get-wx/internal/noaa"
)

func displayForecast(wx Weather) {
	fmt.Printf("\n----------------------------------\n")
	fmt.Println("Last Updated:", wx.Updated.Local().Format("2006-01-02 15:04 MST"))
	fmt.Printf("Elevation: %.2f m, %.2f ft MSL\n", wx.Elevation, wx.ElevationFT)
	fmt.Printf("Forecast temperature is: %.1f° C, %.1f° F\n", wx.Temperature, wx.TemperatureF)
	fmt.Printf("Forecast air pressure is %.2f hPa, %.2f inHg\n", wx.Pressure, wx.PressureInHg)
	fmt.Printf("Forecast wind speed is: %.1f m/s, %.1f MpH, %.1f KTS\n", wx.WindSpeed, wx.WindSpeedMPH, wx.WindSpeedKTS)
	fmt.Printf("Forecast wind direction is: %03d°\n", wx.WindDirection)
	fmt.Printf("Forecast humidity is: %d%%\n", wx.Humidity)
	fmt.Printf("Forecast sky condition is: %d%% cloudy\n\n", wx.SkyCondition)

	fmt.Println("Forecast for the next hour:")
	fmt.Printf("Sky condition: %v\n", wx.ForecastCondition1H)
	fmt.Printf("Preciptitation %.1f\"\n\n", wx.ForecastPrecipitation1H/25.4)

	fmt.Println("Forecast for the next 6 hours:")
	fmt.Printf("Sky condition: %v\n", wx.ForecastCondition6H)
	fmt.Printf("Preciptitation: %.1f\"\n\n", wx.ForecastPrecipitation6H/25.4)

	fmt.Println("Forecast for the next 12 hours:")
	fmt.Printf("Sky condition: %v\n", wx.ForecastCondition12H)
}

func displayMeasurement(label string, value *float64) {
	if value != nil {
		fmt.Println(label, *value)
	} else {
		fmt.Println(label, "No Data")
	}
}

func displayObservations(observations noaa.ObservationResponse) {
	obs := observations.Properties
	fmt.Printf("\n----------------------------------\n")
	fmt.Println("Current NOAA measured observations at", obs.StationID, obs.StationName)
	displayMeasurement("Temperature:", obs.Temperature.Value)
	displayMeasurement("Dewpoint:", obs.Dewpoint.Value)
	displayMeasurement("WindDirection", obs.WindDirection.Value)
	displayMeasurement("WindSpeed", obs.WindSpeed.Value)
	displayMeasurement("WindGust", obs.WindGust.Value)
	displayMeasurement("BarometricPressure", obs.BarometricPressure.Value)
	displayMeasurement("SeaLevelPressure", obs.SeaLevelPressure.Value)
	displayMeasurement("Visibility", obs.Visibility.Value)
	displayMeasurement("MaxTemperatureLast24Hours", obs.MaxTemperatureLast24Hours.Value)
	displayMeasurement("MinTemperatureLast24Hours", obs.MinTemperatureLast24Hours.Value)
	displayMeasurement("PrecipitationLast3Hours", obs.PrecipitationLast3Hours.Value)
	displayMeasurement("RelativeHumidity", obs.RelativeHumidity.Value)
	displayMeasurement("WindChill", obs.WindChill.Value)
	displayMeasurement("HeatIndex", obs.HeatIndex.Value)
}
