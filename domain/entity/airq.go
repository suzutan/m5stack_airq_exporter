package entity

// AirQuality represents air quality measurement data from M5Stack AirQ device
type AirQuality struct {
	// SEN55 sensor data (particle and environmental)
	PM1_0       float64
	PM2_5       float64
	PM4_0       float64
	PM10_0      float64
	Humidity    float64
	Temperature float64
	VOC         int
	NOx         int

	// SCD40 sensor data (CO2 and climate)
	CO2              int
	SCD40Humidity    float64
	SCD40Temperature float64

	// Device info
	Nickname string
}
