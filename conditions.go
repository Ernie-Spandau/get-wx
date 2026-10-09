package main

import (
	"strings"
)

func formatCondition(condition string) string {
	condition = strings.TrimSuffix(condition, "_day")
	condition = strings.TrimSuffix(condition, "_night")
	condition = strings.TrimSuffix(condition, "_polartwilight")
	switch condition {
	case "clearsky":
		return "Clear skies"
	case "fair":
		return "Fair"
	case "partlycloudy":
		return "Partly cloudy"
	case "cloudy":
		return "Cloudy"
	case "rainshowers":
		return "Rain showers"
	case "rainshowersandthunder":
		return "Rain showers and thunder"
	case "sleetshowers":
		return "Sleet showers"
	case "snowshowers":
		return "Snow showers"
	case "rain":
		return "Rain"
	case "heavyrain":
		return "Heavy rain"
	case "heavyrainandthunder":
		return "Heavy rain and thunder"
	case "sleet":
		return "Sleet"
	case "snow":
		return "Snow"
	case "snowandthunder":
		return "Snow and thunder"
	case "fog":
		return "Fog"
	case "sleetshowersandthunder":
		return "Sleet showers and thunder"
	case "snowshowersandthunder":
		return "Snow showers and thunder"
	case "rainandthunder":
		return "Rain and thunder"
	case "sleetandthunder":
		return "Sleet and thunder"
	case "lightrainshowersandthunder":
		return "Light rain showers and thunder"
	case "heavyrainshowersandthunder":
		return "Heavy rain showers and thunder"
	case "lightssleetshowersandthunder":
		return "Light sleet showers and thunder"
	case "heavysleetshowersandthunder":
		return "Heavy sleet showers and thunder"
	case "lightssnowshowersandthunder":
		return "Light snow showers and thunder"
	case "heavysnowshowersandthunder":
		return "Heavy snow showers and thunder"
	case "lightrainandthunder":
		return "Light rain and thunder"
	case "lightsleetandthunder":
		return "Light sleet and thunder"
	case "heavysleetandthunder":
		return "Heavy sleet and thunder"
	case "lightsnowandthunder":
		return "Light snow and thunder"
	case "heavysnowandthunder":
		return "Heavy snow and thunder"
	case "lightrainshowers":
		return "Light rain showers"
	case "heavyrainshowers":
		return "Heavy rain showers"
	case "lightsleetshowers":
		return "Light sleet showers"
	case "heavysleetshowers":
		return "Heavy sleet showers"
	case "lightsnowshowers":
		return "Light snow showers"
	case "heavysnowshowers":
		return "Heavy snow showers"
	case "lightrain":
		return "Light rain"
	case "lightsleet":
		return "Light sleet"
	case "heavysleet":
		return "Heavy sleet"
	case "lightsnow":
		return "Light snow"
	case "heavysnow":
		return "Heavy snow"
	default:
		return "Unknown condition: " + condition
	}
}
