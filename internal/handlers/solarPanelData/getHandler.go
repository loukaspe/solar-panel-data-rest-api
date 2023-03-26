package solarPanelData

import (
	"encoding/csv"
	"github.com/gorilla/mux"
	"github.com/loukaspe/solar-panel-data-crud/internal/core/services"
	apierrors "github.com/loukaspe/solar-panel-data-crud/pkg/errors"
	"github.com/loukaspe/solar-panel-data-crud/pkg/helper"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type GetSolarPanelDataHandler struct {
	SolarPanelDataService        services.SolarPanelDataServiceInterface
	SolarPanelDataEventExtractor helper.SolarPanelDataEventExtractorInterface
	logger                       *log.Logger
}

func NewGetSolarPanelDataHandler(
	service *services.SolarPanelDataService,
	extractor *helper.SolarPanelDataEventExtractor,
	logger *log.Logger,
) *GetSolarPanelDataHandler {
	return &GetSolarPanelDataHandler{
		SolarPanelDataService:        service,
		SolarPanelDataEventExtractor: extractor,
		logger:                       logger,
	}
}

func (handler *GetSolarPanelDataHandler) GetSolarPanelDataController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/csv")

	csvWriter := csv.NewWriter(w)

	response := &GetSolarPanelDataResponse{}

	dataUuid := mux.Vars(r)["id"]
	if dataUuid == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := csvWriter.WriteAll([][]string{{"missing solarPanelData id"}})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in getting solar panel data")

			return
		}

		return
	}

	solarPanelData, err := handler.SolarPanelDataService.GetSolarPanelData(dataUuid)
	if dataNotFoundErrorWrapper, ok := err.(*apierrors.DataNotFoundErrorWrapper); ok {
		handler.logger.WithFields(log.Fields{
			"errorMessage": dataNotFoundErrorWrapper.Unwrap().Error(),
		}).Debug("Error in getting solar panel data")

		w.WriteHeader(dataNotFoundErrorWrapper.ReturnedStatusCode)

		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		handler.logger.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Error("Error in getting solar panel data")

		return
	}

	response.SolarPanelDataEvents, err = handler.SolarPanelDataEventExtractor.ExtractEventsPerParameterIdToCsvForm(
		solarPanelData,
	)

	if malformedEventDataError, ok := err.(apierrors.MalformedEventDataError); ok {
		handler.logger.WithFields(log.Fields{
			"errorMessage": malformedEventDataError.Unwrap().Error(),
		}).Debug("Error in getting solar panel data")

		w.WriteHeader(malformedEventDataError.ReturnedStatusCode)

		err = csvWriter.WriteAll([][]string{{malformedEventDataError.Error()}})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			handler.logger.WithFields(log.Fields{
				"errorMessage": err.Error(),
			}).Error("Error in getting solar panel data")

			return
		}

		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		handler.logger.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Error("Error in getting solar panel data")

		return
	}

	err = csvWriter.WriteAll(response.SolarPanelDataEvents)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		handler.logger.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Error("Error in getting solar panel data")

		return
	}
}
