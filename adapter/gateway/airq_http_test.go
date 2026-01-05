package gateway

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAirQHTTPGateway_Fetch_Success(t *testing.T) {
	responseJSON := `{
		"code": 200,
		"msg": "OK",
		"data": {
			"dataToken": "test-token",
			"dataType": "string",
			"name": "raw",
			"value": "{\"sen55\":{\"pm1.0\":1.5,\"pm2.5\":2.5,\"pm4.0\":4.0,\"pm10.0\":10.0,\"humidity\":32.54,\"temperature\":23.42,\"voc\":75,\"nox\":1},\"scd40\":{\"co2\":725,\"humidity\":17.99,\"temperature\":31.01},\"rtc\":{\"sleep_interval\":60},\"profile\":{\"nickname\":\"AirQ\"}}",
			"createTime": "1703591914",
			"updateTime": "1767573960"
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	gateway := NewAirQHTTPGateway(server.URL, server.Client())
	data, err := gateway.Fetch(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data.PM1_0 != 1.5 {
		t.Errorf("expected PM1_0 to be 1.5, got %f", data.PM1_0)
	}
	if data.PM2_5 != 2.5 {
		t.Errorf("expected PM2_5 to be 2.5, got %f", data.PM2_5)
	}
	if data.PM4_0 != 4.0 {
		t.Errorf("expected PM4_0 to be 4.0, got %f", data.PM4_0)
	}
	if data.PM10_0 != 10.0 {
		t.Errorf("expected PM10_0 to be 10.0, got %f", data.PM10_0)
	}
	if data.Humidity != 32.54 {
		t.Errorf("expected Humidity to be 32.54, got %f", data.Humidity)
	}
	if data.Temperature != 23.42 {
		t.Errorf("expected Temperature to be 23.42, got %f", data.Temperature)
	}
	if data.VOC != 75 {
		t.Errorf("expected VOC to be 75, got %d", data.VOC)
	}
	if data.NOx != 1 {
		t.Errorf("expected NOx to be 1, got %d", data.NOx)
	}
	if data.CO2 != 725 {
		t.Errorf("expected CO2 to be 725, got %d", data.CO2)
	}
	if data.SCD40Humidity != 17.99 {
		t.Errorf("expected SCD40Humidity to be 17.99, got %f", data.SCD40Humidity)
	}
	if data.SCD40Temperature != 31.01 {
		t.Errorf("expected SCD40Temperature to be 31.01, got %f", data.SCD40Temperature)
	}
	if data.Nickname != "AirQ" {
		t.Errorf("expected Nickname to be AirQ, got %s", data.Nickname)
	}
}

func TestAirQHTTPGateway_Fetch_DoubleEscapedJSON(t *testing.T) {
	// This is the actual format returned by ezdata2.m5stack.com API
	responseJSON := `{"code":200,"msg":"OK","data":{"dataToken":"test-token","dataType":"string","name":"raw","value":"{\\\"sen55\\\":{\\\"pm1.0\\\":1.5,\\\"pm2.5\\\":2.5,\\\"pm4.0\\\":4.0,\\\"pm10.0\\\":10.0,\\\"humidity\\\":32.54,\\\"temperature\\\":23.42,\\\"voc\\\":75,\\\"nox\\\":1},\\\"scd40\\\":{\\\"co2\\\":725,\\\"humidity\\\":17.99,\\\"temperature\\\":31.01},\\\"rtc\\\":{\\\"sleep_interval\\\":60},\\\"profile\\\":{\\\"nickname\\\":\\\"AirQ\\\"}}","createTime":"1703591914","updateTime":"1767573960"}}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	gateway := NewAirQHTTPGateway(server.URL, server.Client())
	data, err := gateway.Fetch(context.Background())

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if data.PM1_0 != 1.5 {
		t.Errorf("expected PM1_0 to be 1.5, got %f", data.PM1_0)
	}
	if data.CO2 != 725 {
		t.Errorf("expected CO2 to be 725, got %d", data.CO2)
	}
	if data.Nickname != "AirQ" {
		t.Errorf("expected Nickname to be AirQ, got %s", data.Nickname)
	}
}

func TestAirQHTTPGateway_Fetch_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	gateway := NewAirQHTTPGateway(server.URL, server.Client())
	_, err := gateway.Fetch(context.Background())

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestAirQHTTPGateway_Fetch_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	gateway := NewAirQHTTPGateway(server.URL, server.Client())
	_, err := gateway.Fetch(context.Background())

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestAirQHTTPGateway_Fetch_APIError(t *testing.T) {
	responseJSON := `{
		"code": 500,
		"msg": "Internal Error",
		"data": null
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	gateway := NewAirQHTTPGateway(server.URL, server.Client())
	_, err := gateway.Fetch(context.Background())

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestAirQHTTPGateway_Fetch_InvalidSensorData(t *testing.T) {
	responseJSON := `{
		"code": 200,
		"msg": "OK",
		"data": {
			"dataToken": "test-token",
			"dataType": "string",
			"name": "raw",
			"value": "invalid json",
			"createTime": "1703591914",
			"updateTime": "1767573960"
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(responseJSON))
	}))
	defer server.Close()

	gateway := NewAirQHTTPGateway(server.URL, server.Client())
	_, err := gateway.Fetch(context.Background())

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestAirQHTTPGateway_Fetch_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// This handler should not be reached
		t.Error("handler should not be called")
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	gateway := NewAirQHTTPGateway(server.URL, server.Client())
	_, err := gateway.Fetch(ctx)

	if err == nil {
		t.Error("expected error, got nil")
	}
}
