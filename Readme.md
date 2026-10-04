# get-wx

A simple command-line weather utility written in Go.

`get-wx` retrieves current weather conditions for either a named location or a set of latitude/longitude coordinates.

This is my first Go project, built as a hands-on exercise for learning the language.

## Features

- Look up weather by location name
- Look up weather by latitude/longitude
- Geocodes location names using OpenStreetMap Nominatim
- Retrieves weather data from the MET Norway Locationforecast API
- Displays:
  - Temperature in Celsius and Fahrenheit
  - Atmospheric pressure in hPa and inHg
  - Wind direction
  - Wind speed
  - Relative humidity
  - Cloud cover

## Usage

### Location

Multi-word locations should be quoted.

    get-wx -l "Fallon, Nevada"

Example output:

    Weather for Fallon, Nevada
    Current temperature is: 24.7°C, 76.5°F
    Current air pressure is 1012.10hPa, 29.89inHg
    Current wind direction is: 44°
    Current wind speed is: 3
    Current humidity is: 21%
    Current sky condition is 0% cloudy

### Coordinates

    get-wx -c 30.7621,-86.5705

Coordinates must be supplied as:

    latitude,longitude

## Building

Requires Go.

Clone the repository and build:

    git clone <repository-url>
    cd get-wx
    go build

Or install it into your Go binary directory:

    go install

Make sure your Go binary directory is in your `PATH`.

For example, with fish:

    fish_add_path ~/go/bin

## Data Sources

Location searches use the OpenStreetMap Nominatim geocoding service.

Weather data is provided by the MET Norway Locationforecast API.

## Status

Early learning project. Expect changes as I continue learning Go.