package services

import (
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/ports"
	"github.com/loukaspe/solar-panel-data-crud/internal/repositories"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	"net/http"
)

type SolarPanelDataServiceInterface interface {
	GetSolarPanelData(string) (*domain.SolarPanelData, error)
	CreateSolarPanelData(*domain.SolarPanelData) (string, error)
	UpdateSolarPanelData(string, *domain.SolarPanelData) error
	DeleteSolarPanelData(string) error
}

func NewSolarPanelDataService(repository *repositories.SolarPanelDataRepository) *SolarPanelDataService {
	return &SolarPanelDataService{repository: repository}
}

type SolarPanelDataService struct {
	repository ports.SolarPanelDataRepositoryInterface
}

func (service SolarPanelDataService) GetSolarPanelData(uuid string) (*domain.SolarPanelData, error) {
	return service.repository.GetSolarPanelData(uuid)
}

func (service SolarPanelDataService) CreateSolarPanelData(solarPanelData *domain.SolarPanelData) (string, error) {
	if solarPanelData.Solar == nil {
		return "", apierrors.EmptySolarDataError{
			ReturnedStatusCode: http.StatusBadRequest,
		}
	}

	return service.repository.CreateSolarPanelData(solarPanelData)
}

func (service SolarPanelDataService) UpdateSolarPanelData(uuid string, solarPanelData *domain.SolarPanelData) error {
	if solarPanelData.Solar == nil {
		return apierrors.EmptySolarDataError{
			ReturnedStatusCode: http.StatusBadRequest,
		}
	}

	return service.repository.UpdateSolarPanelData(uuid, solarPanelData)
}

func (service SolarPanelDataService) DeleteSolarPanelData(uuid string) error {
	return service.repository.DeleteSolarPanelData(uuid)
}
