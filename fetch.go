package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/sync/errgroup"
)

const (
	version              = "20251023_1426"
	baseUrl              = "https://www.meteoswiss.admin.ch/product/output/"
	weatherBaseApi       = "weather-widget/forecast"
	forecastChartBaseApi = "forecast-chart"
)

func fetchVersions() (map[string]string, error) {
	v, err := url.Parse(baseUrl + "versions.json")
	if err != nil {
		return nil, fmt.Errorf("invalid api %s: %w", v, err)
	}
	return fetch[map[string]string](v)
}

func fetchAll(code string) (APIResponse, error) {
	apiVersions, err := fetchVersions()
	if err != nil {
		return APIResponse{}, fmt.Errorf("error fetching API versions: %w", err)
	}

	var widgetData MSWidget
	var forecastChartData []MSForecastChart

	var eg errgroup.Group

	eg.Go(func() error {
		weatherApi, err := buildURL(weatherBaseApi, apiVersions, code)
		if err != nil {
			return fmt.Errorf("can't build weather url: %w", err)
		}
		widgetData, err = fetch[MSWidget](weatherApi)
		if err != nil {
			return fmt.Errorf("can't fetch weather: %w", err)
		}
		return nil
	})
	eg.Go(func() error {
		forecastChartApi, err := buildURL(forecastChartBaseApi, apiVersions, code)
		if err != nil {
			return fmt.Errorf("can't build forecastChart url: %w", err)
		}
		forecastChartData, err = fetch[[]MSForecastChart](forecastChartApi)
		if err != nil {
			return fmt.Errorf("can't fetch forecastChart: %w", err)
		}
		return nil
	})
	if err := eg.Wait(); err != nil {
		return APIResponse{}, err
	}
	// TODO: save forecast chart data into the API Response
	_ = forecastChartData

	temp, err := strconv.ParseFloat(widgetData.Data.Current.Temperature, 64)
	if err != nil {
		return APIResponse{}, fmt.Errorf("can't convert temperature to float64: %w", err)
	}
	apiData := APIResponse{
		CityName:           widgetData.Data.CityName,
		Timestamp:          widgetData.Data.Timestamp,
		SymbolID:           widgetData.Data.SymbolID,
		CurrentTemperature: temp,
	}
	return apiData, nil
}

func fetch[T any](url *url.URL) (T, error) {
	var t T
	res, err := http.Get(url.String())
	if err != nil {
		return t, fmt.Errorf("failed to fetch data: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		return t, fmt.Errorf("unexpected status code for %s: %d", url, res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return t, fmt.Errorf("failed to read response body: %w", err)
	}
	if err := json.Unmarshal(body, &t); err != nil {
		var errResponse MSError
		if altErr := json.Unmarshal(body, &errResponse); altErr == nil {
			slog.Error("meteoswiss api returned error", "status", errResponse.Status, "message", errResponse.Message)
			return t, errResponse
		}
		return t, fmt.Errorf("failed to decode response body: %w", err)
	}
	return t, nil
}

func buildURL(api string, apiVersions map[string]string, code string) (*url.URL, error) {
	for len(code) < 6 {
		code += "0"
	}
	version, ok := apiVersions[api]
	if !ok {
		return nil, fmt.Errorf("version not found for api: %s", api)
	}
	u := fmt.Sprintf("%s%s/version__%s/en/%s.json", baseUrl, api, version, code)
	apiURL, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("invalid api %s: %w", u, err)
	}
	query := apiURL.Query()
	query.Set("plz", code)
	apiURL.RawQuery = query.Encode()
	return apiURL, nil
}
