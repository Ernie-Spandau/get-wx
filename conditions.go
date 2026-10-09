package main

func formatCondition(condition string) string {
	switch condition {
	case "clearsky_day", "clearsky_night", "clearsky_polartwilight":
		return "Clear Skies"
	case "partlycloudy_day", "partlycloudy_night", "partlycloudy_polartwilight":
		return "Partly Cloudy"
	case "cloudy":
		return "Cloudy"
	case "lightrain":
		return "Light rain"
	case "rain":
		return "Rain"
	default:
		return condition
	}
}
