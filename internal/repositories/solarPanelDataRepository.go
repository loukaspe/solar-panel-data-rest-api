package repositories

import (
	"errors"
	"github.com/google/uuid"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	"github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	"net/http"
)

type SolarPanelDataRepository struct {
	db SolarPanelDataDB
}

func NewSolarPanelDataRepository(db SolarPanelDataDB) *SolarPanelDataRepository {
	return &SolarPanelDataRepository{db: db}
}

func (repo *SolarPanelDataRepository) CreateSolarPanelData(solarPanelData *domain.SolarPanelData) (string, error) {
	insertedId := uuid.New().String()

	dao := SolarPanelData{
		Solar: solarPanelData.Solar,
		Wind:  solarPanelData.Wind,
	}

	repo.db[insertedId] = &dao

	return insertedId, nil
}

func (repo *SolarPanelDataRepository) GetSolarPanelData(uuid string) (*domain.SolarPanelData, error) {
	var err error

	retrievedSolarPanelData, exists := repo.db[uuid]

	if !exists {
		return &domain.SolarPanelData{},
			&apierrors.DataNotFoundErrorWrapper{
				ReturnedStatusCode: http.StatusNoContent,
				OriginalError:      errors.New("uuid " + uuid + " not found"),
			}
	}

	return &domain.SolarPanelData{
		Solar: retrievedSolarPanelData.Solar,
		Wind:  retrievedSolarPanelData.Wind,
	}, err
}

func (repo *SolarPanelDataRepository) UpdateSolarPanelData(uuid string, solarPanelData *domain.SolarPanelData) error {
	_, exists := repo.db[uuid]

	if !exists {
		return &apierrors.DataNotFoundErrorWrapper{
			ReturnedStatusCode: http.StatusNotFound,
			OriginalError:      errors.New("uuid " + uuid + " not found"),
		}
	}

	dao := SolarPanelData{
		Solar: solarPanelData.Solar,
		Wind:  solarPanelData.Wind,
	}

	repo.db[uuid] = &dao

	return nil
}

func (repo *SolarPanelDataRepository) DeleteSolarPanelData(uuid string) error {
	delete(repo.db, uuid)

	return nil
}
