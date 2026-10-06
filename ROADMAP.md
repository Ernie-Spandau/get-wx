# Roadmap

`get-wx` started as a small command-line client for current MET Norway weather data. The long-term goal is to keep the weather and location logic independent of any one interface so the same core can support CLI, terminal, desktop, and lightweight widget use cases.

The roadmap is intentionally incremental. Each phase should leave the application useful on its own while making the next interface possible without rewriting the weather client.

## v0.2 — Reliability and test coverage

- Add unit tests for coordinate parsing and validation
- Add tests for weather conversions and formatting
- Add representative API-response fixtures
- Improve error handling for malformed input, network failures, HTTP errors, and incomplete API responses
- Keep external-service behavior isolated enough to test without depending on live API calls

## v0.3 — Core package refactor

Separate the reusable weather client from the command-line entry point.

Proposed structure:

```text
get-wx/
├── cmd/
│   └── get-wx/
│       └── main.go
├── weather/
│   ├── location.go
│   ├── met.go
│   ├── forecast.go
│   └── format.go
└── go.mod
```

Goals:

- Move MET Norway communication into a reusable package
- Move geocoding/location handling out of `main`
- Define application-level weather data independently of the MET response schema
- Centralize unit conversions and presentation formatting
- Keep `main` focused on argument parsing, orchestration, and exit behavior

## v0.4 — CLI expansion

Make the CLI useful both interactively and as input to other programs.

Planned additions:

- Cleaner human-readable output
- Consistent wind heading formatting
- Machine-readable JSON output
- Selectable single-field output, for example temperature, wind, or pressure
- Predictable stdout/stderr behavior for scripting
- Useful exit codes
- Optional cached results to avoid unnecessary API requests

Example targets:

```sh
get-wx -l "Crestview, FL"
get-wx -c 30.7621,-86.5705
get-wx --json
get-wx --field temperature
get-wx --field wind
```

## v0.5 — Terminal UI

Add an interactive terminal interface backed by the same weather package.

Planned features:

- Current conditions dashboard
- Manual refresh
- Forecast view
- Keyboard navigation
- Location selection
- Non-blocking weather updates
- Clear loading and error states

The TUI should remain a frontend; API access, conversions, caching, and forecast interpretation should stay in the core package.

## v0.6 — Desktop GUI

Add a small graphical desktop client without replacing the CLI or TUI.

Planned features:

- Current conditions
- Short forecast
- Location management
- Refresh controls
- Settings for units and presentation
- Cross-platform builds where practical

The GUI should consume the same application-level weather data as the other interfaces.

## v0.7 — Desktop and Conky integration

Make `get-wx` useful as an always-available desktop weather source.

Planned options:

- Conky-friendly output
- Single-value queries suitable for shell/desktop integration
- Cached hourly refresh
- A lightweight frameless desktop weather display
- Transparent or minimal presentation suitable for placement on the desktop background
- Configurable refresh interval and location

The first implementation can use the CLI as a data source for Conky before introducing a dedicated desktop widget.

## Later

Possible future work after the core interfaces are stable:

- Multi-day forecasts
- ICAO location lookup as an additional location source
  - Resolve ICAO identifiers to airport/facility names and coordinates
  - Feed resolved coordinates through the existing weather pipeline
  - Support ICAO input across CLI and future TUI/GUI interfaces
- NOTAM integration after ICAO/location handling is established
  - Query authoritative NOTAM data by ICAO location
  - Explore geographic NOTAM queries by coordinates and radius
  - Keep NOTAM retrieval/presentation separate from the core weather model
- Multiple saved locations
- Automatic approximate location via IP geolocation when neither `-l` nor `-c` is supplied
  - Return latitude/longitude (and a useful location label where available) directly from the geolocation service
  - Treat IP-derived location as approximate; explicit `-l` and `-c` input always takes precedence
  - Keep the provider behind the location layer so it can be changed without affecting the weather client
- Config file support
- Additional weather details and derived values
- Release builds for Linux, Windows, and macOS
- Automated cross-platform builds and GitHub releases
- Exploration of mobile reuse where the shared Go weather package is a good fit

## Design direction

A central goal is to avoid building separate weather applications for each interface. The intended architecture is:

```text
                 ┌── CLI
                 │
MET / location ──┼── TUI
      │          │
  weather core ──┼── GUI
                 │
                 └── Conky / desktop widget
```

Network access, geocoding, weather models, conversions, validation, and caching belong in reusable code. Each interface should be responsible primarily for input and presentation.
