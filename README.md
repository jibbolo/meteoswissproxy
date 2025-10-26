# MeteoSwiss Proxy

A lightweight Go-based HTTP proxy server for fetching weather data from MeteoSwiss API with built-in caching.

## Features

- **Simple HTTP API**: Single endpoint to fetch weather data by location code
- **Response Caching**: Built-in 5-minute cache to reduce API calls
- **Request Validation**: Middleware for validating location codes
- **Concurrent Fetching**: Efficient parallel data fetching
- **Structured Logging**: Built-in logging with `slog`
- **Configurable Port**: Set via environment variable

## Installation

```bash
go get github.com/jibbolo/meteoswiss
```

## Usage

### Running the Server

```bash
go run .
```

Or with a custom port:

```bash
PORT=3000 go run .
```

### API Endpoint

```
GET /{code}
```

**Parameters:**
- `code` (path parameter): The MeteoSwiss location code

**Example Request:**

```bash
curl http://localhost:8080/8200
```

**Example Response:**

```json
{
  "city_name": "Zürich",
  "current_temperature": 15.5,
  "timestamp": 1234567890,
  "symbol_id": 1
}
```

## Configuration

- `PORT`: Server port (default: `8080`)

## Building

```bash
go build -o meteoswissproxy
```

## Requirements

- Go 1.25.3 or higher

## Dependencies

- `golang.org/x/sync` - For concurrent operations

## License

See [LICENSE](LICENSE) file for details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for version history.