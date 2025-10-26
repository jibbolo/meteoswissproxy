package main

import "fmt"

type APIResponse struct {
	CityName           string  `json:"city_name"`
	CurrentTemperature float64 `json:"current_temperature"`
	Timestamp          int     `json:"timestamp"`
	SymbolID           int     `json:"symbol_id"`
}

type MSError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e MSError) Error() string {
	return fmt.Sprintf("API error: %d %s", e.Status, e.Message)
}

type MSWidget struct {
	Data MSWidgetData `json:"data"`
}

type MSWidgetData struct {
	Altitude   int           `json:"altitude"`
	CityName   string        `json:"city_name"`
	CoordX     int           `json:"coord_x"`
	CoordY     int           `json:"coord_y"`
	Current    MSCurrent     `json:"current"`
	Forecasts  []MSForecasts `json:"forecasts"`
	LocationID string        `json:"location_id"`
	Timestamp  int           `json:"timestamp"`
	SymbolID   int           `json:"weather_symbol_id"`
}
type MSCurrent struct {
	Temperature    string `json:"temperature"`
	SymbolIDString string `json:"weather_symbol_id"`
}
type MSForecasts struct {
	PrecipMax      string `json:"precip_max"`
	Noon           int64  `json:"noon"`
	PrecipMin      string `json:"precip_min"`
	TempHigh       string `json:"temp_high"`
	Weekday        string `json:"weekday"`
	TempLow        string `json:"temp_low"`
	PrecipMean     string `json:"precip_mean"`
	SymbolIDString string `json:"weather_symbol_id"`
}

type MSForecastChart struct {
	CurrentTime           int64                          `json:"current_time"`
	CurrentTimeString     string                         `json:"current_time_string"`
	DayString             string                         `json:"day_string"`
	MaxDate               int64                          `json:"max_date"`
	MinDate               int64                          `json:"min_date"`
	NewDay                int64                          `json:"new_day,omitempty"`
	Rainfall              [][]float64                    `json:"rainfall"`
	Sunrise               int64                          `json:"sunrise"`
	Sunset                int64                          `json:"sunset"`
	Sunshine              [][]float64                    `json:"sunshine"`
	SymbolDay             MSForecastChartWeatherSymbol   `json:"symbol_day"`
	Symbols               []MSForecastChartWeatherSymbol `json:"symbols"`
	Temperature           [][]float64                    `json:"temperature"`
	VarianceRain          [][]float64                    `json:"variance_rain"`
	VarianceRange         [][]float64                    `json:"variance_range"`
	Wind                  MSForecastChartWind            `json:"wind"`
	WindGustPeak          MSForecastChartWindGustPeak    `json:"wind_gust_peak"`
	WindGustSpeedVariance [][]float64                    `json:"wind_gust_speed_variance"`
	WindSpeedVariance     [][]float64                    `json:"wind_speed_variance"`
}

type MSForecastChartWindGustPeak struct {
	Data [][]float64 `json:"data"`
}
type MSForecastChartWeatherSymbol struct {
	Timestamp int64 `json:"timestamp"`
	SymbolID  int   `json:"weather_symbol_id"`
}

type MSForecastChartWind struct {
	Data    [][]float64                 `json:"data"`
	Symbols []MSForecastChartWindSymbol `json:"symbols"`
}
type MSForecastChartWindSymbol struct {
	SymbolIDString string `json:"symbol_id"`
	Timestamp      int64  `json:"timestamp"`
}
