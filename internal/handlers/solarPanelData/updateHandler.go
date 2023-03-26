package solarPanelData

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/domain"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/services"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type UpdateSolarPanelDataHandler struct {
	SolarPanelDataService services.SolarPanelDataServiceInterface
	logger                *log.Logger
}

func NewUpdateSolarPanelDataHandler(
	service *services.SolarPanelDataService,
	logger *log.Logger,
) *UpdateSolarPanelDataHandler {
	return &UpdateSolarPanelDataHandler{
		SolarPanelDataService: service,
		logger:                logger,
	}
}

func (handler *UpdateSolarPanelDataHandler) UpdateSolarPanelDataController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := &UpdateSolarPanelDataResponse{}
	solarPanelDataRequest := &Dto{}

	err := json.NewDecoder(r.Body).Decode(solarPanelDataRequest)
	if err != nil {
		handler.logger.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Error("Error in updating solar panel data")

		w.WriteHeader(http.StatusBadRequest)
		response.ErrorMessage = "malformed solar panel data request"
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in updating solar panel data")

			return
		}

		return
	}

	domainSolarPanelData := &domain.SolarPanelData{
		Solar: solarPanelDataRequest.Solar,
		Wind:  solarPanelDataRequest.Wind,
	}

	uuid := mux.Vars(r)["id"]
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.ErrorMessage = "missing solarPanelData id"
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in updating solar panel data")

			return
		}

		return
	}

	err = handler.SolarPanelDataService.UpdateSolarPanelData(uuid, domainSolarPanelData)
	if emptySolarDataError, ok := err.(apierrors.EmptySolarDataError); ok {
		w.WriteHeader(emptySolarDataError.ReturnedStatusCode)

		response.ErrorMessage = err.Error()
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in updating solar panel data")

			return
		}

		return
	}

	if dataNotFoundErrorWrapper, ok := err.(*apierrors.DataNotFoundErrorWrapper); ok {
		handler.logger.WithFields(log.Fields{
			"errorMessage": dataNotFoundErrorWrapper.Unwrap().Error(),
		}).Debug("Error in updating solar panel data")

		w.WriteHeader(dataNotFoundErrorWrapper.ReturnedStatusCode)

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
			}).Error("Error in updating solar panel data")

			return
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}
