package main

import (
	"fmt"
)

func displayWeather(wx Weather) {
	fmt.Printf("Elevation: %.2f m, %.2f ft MSL\n", wx.Elevation, wx.ElevationFT)
	fmt.Printf("Current temperature is: %.1f° C, %.1f° F\n", wx.Temperature, wx.TemperatureF)
	fmt.Printf("Current air pressure is %.2f hPa, %.2f inHg\n", wx.Pressure, wx.PressureInHg)
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
