package gateway

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/suzutan/m5stack_airq_exporter/domain/entity"
)

// PrometheusMetricsGateway implements MetricsRepository using Prometheus client
type PrometheusMetricsGateway struct {
	// SEN55 sensor metrics
	pm1_0       prometheus.Gauge
	pm2_5       prometheus.Gauge
	pm4_0       prometheus.Gauge
	pm10_0      prometheus.Gauge
	humidity    prometheus.Gauge
	temperature prometheus.Gauge
	voc         prometheus.Gauge
	nox         prometheus.Gauge

	// SCD40 sensor metrics
	co2              prometheus.Gauge
	scd40Humidity    prometheus.Gauge
	scd40Temperature prometheus.Gauge
}

// NewPrometheusMetricsGateway creates a new PrometheusMetricsGateway and registers metrics
func NewPrometheusMetricsGateway(registry prometheus.Registerer) *PrometheusMetricsGateway {
	g := &PrometheusMetricsGateway{
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

	// Register all metrics
	registry.MustRegister(
		g.pm1_0,
		g.pm2_5,
		g.pm4_0,
		g.pm10_0,
		g.humidity,
		g.temperature,
		g.voc,
		g.nox,
		g.co2,
		g.scd40Humidity,
		g.scd40Temperature,
	)

	return g
}

// Update updates the Prometheus metrics with the given air quality data
func (g *PrometheusMetricsGateway) Update(data *entity.AirQuality) {
	g.pm1_0.Set(data.PM1_0)
	g.pm2_5.Set(data.PM2_5)
	g.pm4_0.Set(data.PM4_0)
	g.pm10_0.Set(data.PM10_0)
	g.humidity.Set(data.Humidity)
	g.temperature.Set(data.Temperature)
	g.voc.Set(float64(data.VOC))
	g.nox.Set(float64(data.NOx))
	g.co2.Set(float64(data.CO2))
	g.scd40Humidity.Set(data.SCD40Humidity)
	g.scd40Temperature.Set(data.SCD40Temperature)
}
