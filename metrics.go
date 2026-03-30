package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics holds Prometheus gauge metrics for air quality measurements
type Metrics struct {
	pm1_0       prometheus.Gauge
	pm2_5       prometheus.Gauge
	pm4_0       prometheus.Gauge
	pm10_0      prometheus.Gauge
	humidity    prometheus.Gauge
	temperature prometheus.Gauge
	voc         prometheus.Gauge
	nox         prometheus.Gauge

	co2              prometheus.Gauge
	scd40Humidity    prometheus.Gauge
	scd40Temperature prometheus.Gauge
}

// NewMetrics creates and registers Prometheus gauge metrics
func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		pm1_0: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_pm1_0",
			Help: "PM1.0 concentration in µg/m³",
		}),
		pm2_5: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_pm2_5",
			Help: "PM2.5 concentration in µg/m³",
		}),
		pm4_0: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_pm4_0",
			Help: "PM4.0 concentration in µg/m³",
		}),
		pm10_0: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_pm10_0",
			Help: "PM10.0 concentration in µg/m³",
		}),
		humidity: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_humidity",
			Help: "Relative humidity in % (SEN55)",
		}),
		temperature: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_temperature",
			Help: "Temperature in °C (SEN55)",
		}),
		voc: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_voc",
			Help: "VOC index",
		}),
		nox: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_nox",
			Help: "NOx index",
		}),
		co2: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_co2",
			Help: "CO2 concentration in ppm",
		}),
		scd40Humidity: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_scd40_humidity",
			Help: "Relative humidity in % (SCD40)",
		}),
		scd40Temperature: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "airq_scd40_temperature",
			Help: "Temperature in °C (SCD40)",
		}),
	}

	reg.MustRegister(
		m.pm1_0,
		m.pm2_5,
		m.pm4_0,
		m.pm10_0,
		m.humidity,
		m.temperature,
		m.voc,
		m.nox,
		m.co2,
		m.scd40Humidity,
		m.scd40Temperature,
	)

	return m
}

// Update updates the Prometheus metrics with the given air quality data
func (m *Metrics) Update(data *AirQuality) {
	m.pm1_0.Set(data.PM1_0)
	m.pm2_5.Set(data.PM2_5)
	m.pm4_0.Set(data.PM4_0)
	m.pm10_0.Set(data.PM10_0)
	m.humidity.Set(data.Humidity)
	m.temperature.Set(data.Temperature)
	m.voc.Set(float64(data.VOC))
	m.nox.Set(float64(data.NOx))
	m.co2.Set(float64(data.CO2))
	m.scd40Humidity.Set(data.SCD40Humidity)
	m.scd40Temperature.Set(data.SCD40Temperature)
}
