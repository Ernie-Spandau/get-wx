package main

import (
	"flag"
	"fmt"
	"get-wx/internal/met"
	"get-wx/internal/noaa"
)

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

		fmt.Printf("Forecast weather for %s (%.4f, %.4f)\n", *location, latitude, longitude)

	} else {

		latitude, longitude, err = parseCoordinates(*coordinates)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Forecast weather for %.4f, %.4f\n", latitude, longitude)
	}
	elevation, err := getElevation(latitude, longitude)
	if err != nil {
		fmt.Println(err)
		return
	}

	forecast, err := met.GetForecast(latitude, longitude, elevation)
	if err != nil {
		fmt.Println("MET Norway", err)
		return
	}

	stations, err := noaa.GetStations(latitude, longitude)
	if err != nil {
		fmt.Println("NOAA", err)
		return
	}

	observations, err := noaa.GetStationWx(stations)
	if err != nil {
		fmt.Println("NOAA:", err)
		return
	}
	wx := convertWeather(forecast, elevation)
	displayObservations(observations)
	displayForecast(wx)
}
