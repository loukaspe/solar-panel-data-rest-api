package solarPanelData

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/services"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type DeleteSolarPanelDataHandler struct {
	SolarPanelDataService services.SolarPanelDataServiceInterface
	logger                *log.Logger
}

func NewDeleteSolarPanelDataHandler(
	service *services.SolarPanelDataService,
	logger *log.Logger,
) *DeleteSolarPanelDataHandler {
	return &DeleteSolarPanelDataHandler{
		SolarPanelDataService: service,
		logger:                logger,
	}
}

func (handler *DeleteSolarPanelDataHandler) DeleteSolarPanelDataController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var err error
	response := &DeleteSolarPanelDataResponse{}

	uuid := mux.Vars(r)["id"]
	if uuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		response.ErrorMessage = "missing solarPanelData id"
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in deleting solar panel data")

			return
		}

		return
	}

	err = handler.SolarPanelDataService.DeleteSolarPanelData(uuid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response.ErrorMessage = err.Error()
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in deleting solar panel data")

			return
		}

		return
	}

	w.WriteHeader(http.StatusOK)
}
