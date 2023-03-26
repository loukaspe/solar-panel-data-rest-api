package ports

import "github.com/loukaspe/solar-panel-data-crud/internal/core/domain"

type SolarPanelDataRepositoryInterface interface {
	GetSolarPanelData(uuid string) (*domain.SolarPanelData, error)
	CreateSolarPanelData(*domain.SolarPanelData) (string, error)
	UpdateSolarPanelData(string, *domain.SolarPanelData) error
	DeleteSolarPanelData(string) error
}
