package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/suzutan/m5stack_airq_exporter/domain/entity"
)

// HTTPClient interface for HTTP client operations (allows mocking)
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AirQHTTPGateway implements AirQRepository using HTTP client
type AirQHTTPGateway struct {
	url    string
	client HTTPClient
}

// NewAirQHTTPGateway creates a new AirQHTTPGateway with the given URL and HTTP client
func NewAirQHTTPGateway(url string, client HTTPClient) *AirQHTTPGateway {
	return &AirQHTTPGateway{
		url:    url,
		client: client,
	}
}

// apiResponse represents the top-level response from ezdata2.m5stack.com API
type apiResponse struct {
	Code int       `json:"code"`
	Msg  string    `json:"msg"`
	Data *dataBody `json:"data"`
}

// dataBody represents the data field in the API response
type dataBody struct {
	DataToken  string `json:"dataToken"`
	DataType   string `json:"dataType"`
	Name       string `json:"name"`
	Value      string `json:"value"`
	CreateTime string `json:"createTime"`
	UpdateTime string `json:"updateTime"`
}

// sensorData represents the parsed sensor data from the Value field
type sensorData struct {
	SEN55   sen55Data   `json:"sen55"`
	SCD40   scd40Data   `json:"scd40"`
	RTC     rtcData     `json:"rtc"`
	Profile profileData `json:"profile"`
}

// sen55Data represents data from the SEN55 sensor
type sen55Data struct {
	PM1_0       float64 `json:"pm1.0"`
	PM2_5       float64 `json:"pm2.5"`
	PM4_0       float64 `json:"pm4.0"`
	PM10_0      float64 `json:"pm10.0"`
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
	VOC         int     `json:"voc"`
	NOx         int     `json:"nox"`
}

// scd40Data represents data from the SCD40 sensor
type scd40Data struct {
	CO2         int     `json:"co2"`
	Humidity    float64 `json:"humidity"`
	Temperature float64 `json:"temperature"`
}

// rtcData represents real-time clock configuration
type rtcData struct {
	SleepInterval int `json:"sleep_interval"`
}

// profileData represents device profile information
type profileData struct {
	Nickname string `json:"nickname"`
}

// Fetch retrieves the latest air quality data from the API
func (g *AirQHTTPGateway) Fetch(ctx context.Context) (*entity.AirQuality, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, g.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	if apiResp.Code != 200 {
		return nil, fmt.Errorf("API error: code=%d, msg=%s", apiResp.Code, apiResp.Msg)
	}

	if apiResp.Data == nil {
		return nil, fmt.Errorf("API response data is nil")
	}

	// The value field may be double-escaped JSON, try to unescape it
	valueStr := apiResp.Data.Value
	if unquoted, err := strconv.Unquote(`"` + valueStr + `"`); err == nil {
		valueStr = unquoted
	}

	var sensor sensorData
	if err := json.Unmarshal([]byte(valueStr), &sensor); err != nil {
		return nil, fmt.Errorf("failed to parse sensor data: %w", err)
	}

	return &entity.AirQuality{
		PM1_0:            sensor.SEN55.PM1_0,
		PM2_5:            sensor.SEN55.PM2_5,
		PM4_0:            sensor.SEN55.PM4_0,
		PM10_0:           sensor.SEN55.PM10_0,
		Humidity:         sensor.SEN55.Humidity,
		Temperature:      sensor.SEN55.Temperature,
		VOC:              sensor.SEN55.VOC,
		NOx:              sensor.SEN55.NOx,
		CO2:              sensor.SCD40.CO2,
		SCD40Humidity:    sensor.SCD40.Humidity,
		SCD40Temperature: sensor.SCD40.Temperature,
		Nickname:         sensor.Profile.Nickname,
	}, nil
}
