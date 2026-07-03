package main

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMetrics_Update(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	data := &AirQuality{
		PM1_0:            1.5,
		PM2_5:            2.5,
		PM4_0:            4.0,
		PM10_0:           10.0,
		Humidity:         32.54,
		Temperature:      23.42,
		VOC:              75,
		NOx:              1,
		CO2:              725,
		SCD40Humidity:    17.99,
		SCD40Temperature: 31.01,
		Nickname:         "AirQ",
	}

	metrics.Update(data)

	expected := `
		# HELP airq_pm1_0 PM1.0 concentration in µg/m³
		# TYPE airq_pm1_0 gauge
		airq_pm1_0 1.5
	`
	if err := testutil.CollectAndCompare(metrics.pm1_0, strings.NewReader(expected)); err != nil {
		t.Errorf("PM1.0 metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_pm2_5 PM2.5 concentration in µg/m³
		# TYPE airq_pm2_5 gauge
		airq_pm2_5 2.5
	`
	if err := testutil.CollectAndCompare(metrics.pm2_5, strings.NewReader(expected)); err != nil {
		t.Errorf("PM2.5 metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_pm4_0 PM4.0 concentration in µg/m³
		# TYPE airq_pm4_0 gauge
		airq_pm4_0 4
	`
	if err := testutil.CollectAndCompare(metrics.pm4_0, strings.NewReader(expected)); err != nil {
		t.Errorf("PM4.0 metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_pm10_0 PM10.0 concentration in µg/m³
		# TYPE airq_pm10_0 gauge
		airq_pm10_0 10
	`
	if err := testutil.CollectAndCompare(metrics.pm10_0, strings.NewReader(expected)); err != nil {
		t.Errorf("PM10.0 metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_humidity Relative humidity in % (SEN55)
		# TYPE airq_humidity gauge
		airq_humidity 32.54
	`
	if err := testutil.CollectAndCompare(metrics.humidity, strings.NewReader(expected)); err != nil {
		t.Errorf("Humidity metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_temperature Temperature in °C (SEN55)
		# TYPE airq_temperature gauge
		airq_temperature 23.42
	`
	if err := testutil.CollectAndCompare(metrics.temperature, strings.NewReader(expected)); err != nil {
		t.Errorf("Temperature metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_voc VOC index
		# TYPE airq_voc gauge
		airq_voc 75
	`
	if err := testutil.CollectAndCompare(metrics.voc, strings.NewReader(expected)); err != nil {
		t.Errorf("VOC metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_nox NOx index
		# TYPE airq_nox gauge
		airq_nox 1
	`
	if err := testutil.CollectAndCompare(metrics.nox, strings.NewReader(expected)); err != nil {
		t.Errorf("NOx metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_co2 CO2 concentration in ppm
		# TYPE airq_co2 gauge
		airq_co2 725
	`
	if err := testutil.CollectAndCompare(metrics.co2, strings.NewReader(expected)); err != nil {
		t.Errorf("CO2 metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_scd40_humidity Relative humidity in % (SCD40)
		# TYPE airq_scd40_humidity gauge
		airq_scd40_humidity 17.99
	`
	if err := testutil.CollectAndCompare(metrics.scd40Humidity, strings.NewReader(expected)); err != nil {
		t.Errorf("SCD40 Humidity metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_scd40_temperature Temperature in °C (SCD40)
		# TYPE airq_scd40_temperature gauge
		airq_scd40_temperature 31.01
	`
	if err := testutil.CollectAndCompare(metrics.scd40Temperature, strings.NewReader(expected)); err != nil {
		t.Errorf("SCD40 Temperature metric mismatch: %v", err)
	}
}

func TestMetrics_UpdateMultipleTimes(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	// First update
	data1 := &AirQuality{
		PM2_5: 10.0,
		CO2:   500,
	}
	metrics.Update(data1)

	// Second update with different values
	data2 := &AirQuality{
		PM2_5: 25.0,
		CO2:   800,
	}
	metrics.Update(data2)

	// Should have the latest values
	expected := `
		# HELP airq_pm2_5 PM2.5 concentration in µg/m³
		# TYPE airq_pm2_5 gauge
		airq_pm2_5 25
	`
	if err := testutil.CollectAndCompare(metrics.pm2_5, strings.NewReader(expected)); err != nil {
		t.Errorf("PM2.5 metric should be updated: %v", err)
	}

	expected = `
		# HELP airq_co2 CO2 concentration in ppm
		# TYPE airq_co2 gauge
		airq_co2 800
	`
	if err := testutil.CollectAndCompare(metrics.co2, strings.NewReader(expected)); err != nil {
		t.Errorf("CO2 metric should be updated: %v", err)
	}
}

func TestMetrics_RecordScrapeError(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordScrapeError()

	expected := `
		# HELP airq_scrape_success Whether the last fetch of air quality data succeeded (1) or failed (0)
		# TYPE airq_scrape_success gauge
		airq_scrape_success 0
	`
	if err := testutil.CollectAndCompare(metrics.scrapeSuccess, strings.NewReader(expected)); err != nil {
		t.Errorf("Scrape success metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_scrape_errors_total Total number of failed fetches of air quality data
		# TYPE airq_scrape_errors_total counter
		airq_scrape_errors_total 1
	`
	if err := testutil.CollectAndCompare(metrics.scrapeErrors, strings.NewReader(expected)); err != nil {
		t.Errorf("Scrape errors metric mismatch: %v", err)
	}
}

func TestMetrics_RecordScrapeSuccess(t *testing.T) {
	registry := prometheus.NewRegistry()
	metrics := NewMetrics(registry)

	metrics.RecordScrapeSuccess()

	expected := `
		# HELP airq_scrape_success Whether the last fetch of air quality data succeeded (1) or failed (0)
		# TYPE airq_scrape_success gauge
		airq_scrape_success 1
	`
	if err := testutil.CollectAndCompare(metrics.scrapeSuccess, strings.NewReader(expected)); err != nil {
		t.Errorf("Scrape success metric mismatch: %v", err)
	}

	if ts := testutil.ToFloat64(metrics.lastScrapeSuccessTimestamp); ts <= 0 {
		t.Errorf("Last scrape success timestamp should be greater than 0, got %v", ts)
	}
}
