package solarPanelData

import (
	"encoding/json"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/services"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type CreateSolarPanelDataHandler struct {
	SolarPanelDataService services.SolarPanelDataServiceInterface
	logger                *log.Logger
}

func NewCreateSolarPanelDataHandler(
	service *services.SolarPanelDataService,
	logger *log.Logger,
) *CreateSolarPanelDataHandler {
	return &CreateSolarPanelDataHandler{
		SolarPanelDataService: service,
		logger:                logger,
	}
}

func (handler *CreateSolarPanelDataHandler) CreateSolarPanelDataController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := &CreateSolarPanelDataResponse{}
	solarPanelDataRequest := &Dto{}

	err := json.NewDecoder(r.Body).Decode(solarPanelDataRequest)
	if err != nil {
		handler.logger.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Error("Error in creating solar panel data")

		w.WriteHeader(http.StatusBadRequest)
		response.ErrorMessage = "malformed solar panel data request"
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in creating solar panel data")

			return
		}

		return
	}

	domainSolarPanelData := &domain.SolarPanelData{
		Solar: solarPanelDataRequest.Solar,
		Wind:  solarPanelDataRequest.Wind,
	}

	insertedId, err := handler.SolarPanelDataService.CreateSolarPanelData(domainSolarPanelData)
	if emptySolarDataError, ok := err.(apierrors.EmptySolarDataError); ok {
		w.WriteHeader(emptySolarDataError.ReturnedStatusCode)

		response.ErrorMessage = err.Error()
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in creating solar panel data")

			return
		}

		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.ErrorMessage = err.Error()
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in creating solar panel data")

			return
		}

		return
	}

	w.WriteHeader(http.StatusCreated)

	response.DataSubmitted = solarPanelDataRequest
	response.InsertedId = insertedId

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		handler.logger.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Error("Error in creating solar panel data")

		return
	}
}
