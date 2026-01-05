package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/suzutan/m5stack_airq_exporter/domain/entity"
)

// mockAirQRepository is a mock implementation of AirQRepository for testing
type mockAirQRepository struct {
	data *entity.AirQuality
	err  error
}

func (m *mockAirQRepository) Fetch(ctx context.Context) (*entity.AirQuality, error) {
	return m.data, m.err
}

// mockMetricsRepository is a mock implementation of MetricsRepository for testing
type mockMetricsRepository struct {
	updatedData *entity.AirQuality
	updateCount int
}

func (m *mockMetricsRepository) Update(data *entity.AirQuality) {
	m.updatedData = data
	m.updateCount++
}

func TestFetchAirQUsecase_Execute_Success(t *testing.T) {
	expectedData := &entity.AirQuality{
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

	airqRepo := &mockAirQRepository{data: expectedData}
	metricsRepo := &mockMetricsRepository{}

	usecase := NewFetchAirQUsecase(airqRepo, metricsRepo)
	err := usecase.Execute(context.Background())

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if metricsRepo.updateCount != 1 {
		t.Errorf("expected Update to be called once, got %d", metricsRepo.updateCount)
	}

	if metricsRepo.updatedData != expectedData {
		t.Errorf("expected updated data to match, got %+v", metricsRepo.updatedData)
	}
}

func TestFetchAirQUsecase_Execute_FetchError(t *testing.T) {
	expectedErr := errors.New("fetch error")

	airqRepo := &mockAirQRepository{err: expectedErr}
	metricsRepo := &mockMetricsRepository{}

	usecase := NewFetchAirQUsecase(airqRepo, metricsRepo)
	err := usecase.Execute(context.Background())

	if err == nil {
		t.Error("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}

	if metricsRepo.updateCount != 0 {
		t.Errorf("expected Update not to be called, got %d", metricsRepo.updateCount)
	}
}

func TestFetchAirQUsecase_Execute_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	airqRepo := &mockAirQRepository{err: context.Canceled}
	metricsRepo := &mockMetricsRepository{}

	usecase := NewFetchAirQUsecase(airqRepo, metricsRepo)
	err := usecase.Execute(ctx)

	if err == nil {
		t.Error("expected error, got nil")
	}

	if metricsRepo.updateCount != 0 {
		t.Errorf("expected Update not to be called, got %d", metricsRepo.updateCount)
	}
}
