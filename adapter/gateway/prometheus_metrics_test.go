package gateway

import (
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/suzutan/m5stack_airq_exporter/domain/entity"
)

func TestPrometheusMetricsGateway_Update(t *testing.T) {
	registry := prometheus.NewRegistry()
	gateway := NewPrometheusMetricsGateway(registry)

	data := &entity.AirQuality{
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

	gateway.Update(data)

	// Test PM metrics
	expected := `
		# HELP airq_pm1_0 PM1.0 concentration in µg/m³
		# TYPE airq_pm1_0 gauge
		airq_pm1_0 1.5
	`
	if err := testutil.CollectAndCompare(gateway.pm1_0, strings.NewReader(expected)); err != nil {
		t.Errorf("PM1.0 metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_pm2_5 PM2.5 concentration in µg/m³
		# TYPE airq_pm2_5 gauge
		airq_pm2_5 2.5
	`
	if err := testutil.CollectAndCompare(gateway.pm2_5, strings.NewReader(expected)); err != nil {
		t.Errorf("PM2.5 metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_pm4_0 PM4.0 concentration in µg/m³
		# TYPE airq_pm4_0 gauge
		airq_pm4_0 4
	`
	if err := testutil.CollectAndCompare(gateway.pm4_0, strings.NewReader(expected)); err != nil {
		t.Errorf("PM4.0 metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_pm10_0 PM10.0 concentration in µg/m³
		# TYPE airq_pm10_0 gauge
		airq_pm10_0 10
	`
	if err := testutil.CollectAndCompare(gateway.pm10_0, strings.NewReader(expected)); err != nil {
		t.Errorf("PM10.0 metric mismatch: %v", err)
	}

	// Test environmental metrics
	expected = `
		# HELP airq_humidity Relative humidity in % (SEN55)
		# TYPE airq_humidity gauge
		airq_humidity 32.54
	`
	if err := testutil.CollectAndCompare(gateway.humidity, strings.NewReader(expected)); err != nil {
		t.Errorf("Humidity metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_temperature Temperature in °C (SEN55)
		# TYPE airq_temperature gauge
		airq_temperature 23.42
	`
	if err := testutil.CollectAndCompare(gateway.temperature, strings.NewReader(expected)); err != nil {
		t.Errorf("Temperature metric mismatch: %v", err)
	}

	// Test VOC and NOx
	expected = `
		# HELP airq_voc VOC index
		# TYPE airq_voc gauge
		airq_voc 75
	`
	if err := testutil.CollectAndCompare(gateway.voc, strings.NewReader(expected)); err != nil {
		t.Errorf("VOC metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_nox NOx index
		# TYPE airq_nox gauge
		airq_nox 1
	`
	if err := testutil.CollectAndCompare(gateway.nox, strings.NewReader(expected)); err != nil {
		t.Errorf("NOx metric mismatch: %v", err)
	}

	// Test SCD40 metrics
	expected = `
		# HELP airq_co2 CO2 concentration in ppm
		# TYPE airq_co2 gauge
		airq_co2 725
	`
	if err := testutil.CollectAndCompare(gateway.co2, strings.NewReader(expected)); err != nil {
		t.Errorf("CO2 metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_scd40_humidity Relative humidity in % (SCD40)
		# TYPE airq_scd40_humidity gauge
		airq_scd40_humidity 17.99
	`
	if err := testutil.CollectAndCompare(gateway.scd40Humidity, strings.NewReader(expected)); err != nil {
		t.Errorf("SCD40 Humidity metric mismatch: %v", err)
	}

	expected = `
		# HELP airq_scd40_temperature Temperature in °C (SCD40)
		# TYPE airq_scd40_temperature gauge
		airq_scd40_temperature 31.01
	`
	if err := testutil.CollectAndCompare(gateway.scd40Temperature, strings.NewReader(expected)); err != nil {
		t.Errorf("SCD40 Temperature metric mismatch: %v", err)
	}
}

func TestPrometheusMetricsGateway_UpdateMultipleTimes(t *testing.T) {
	registry := prometheus.NewRegistry()
	gateway := NewPrometheusMetricsGateway(registry)

	// First update
	data1 := &entity.AirQuality{
		PM2_5: 10.0,
		CO2:   500,
	}
	gateway.Update(data1)

	// Second update with different values
	data2 := &entity.AirQuality{
		PM2_5: 25.0,
		CO2:   800,
	}
	gateway.Update(data2)

	// Should have the latest values
	expected := `
		# HELP airq_pm2_5 PM2.5 concentration in µg/m³
		# TYPE airq_pm2_5 gauge
		airq_pm2_5 25
	`
	if err := testutil.CollectAndCompare(gateway.pm2_5, strings.NewReader(expected)); err != nil {
		t.Errorf("PM2.5 metric should be updated: %v", err)
	}

	expected = `
		# HELP airq_co2 CO2 concentration in ppm
		# TYPE airq_co2 gauge
		airq_co2 800
	`
	if err := testutil.CollectAndCompare(gateway.co2, strings.NewReader(expected)); err != nil {
		t.Errorf("CO2 metric should be updated: %v", err)
	}
}
