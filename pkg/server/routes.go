package server

import (
	"github.com/loukaspe/solar-panel-data-crud/internal/core/services"
	"github.com/loukaspe/solar-panel-data-crud/internal/handlers"
	"github.com/loukaspe/solar-panel-data-crud/internal/handlers/solarPanelData"
	"github.com/loukaspe/solar-panel-data-crud/internal/repositories"
	"github.com/loukaspe/solar-panel-data-crud/pkg/helper"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (s *Server) initializeRoutes(
	db repositories.SolarPanelDataDB,
	logger *log.Logger,
) {
	// health check
	healthCheckHandler := handlers.NewHealthCheckHandler(logger)
	s.router.HandleFunc("/health-check", healthCheckHandler.HealthCheckController).Methods("GET")

	// solarPanelData
	solarPanelDataRepository := repositories.NewSolarPanelDataRepository(db)
	solarPanelDataService := services.NewSolarPanelDataService(solarPanelDataRepository)
	solarPanelDataEventExtractor := helper.NewSolarPanelDataEventExtractor()

	getSolarPanelDataHandler := solarPanelData.NewGetSolarPanelDataHandler(
		solarPanelDataService,
		solarPanelDataEventExtractor,
		logger,
	)
	createSolarPanelDataHandler := solarPanelData.NewCreateSolarPanelDataHandler(
		solarPanelDataService,
		logger,
	)
	deleteSolarPanelDataHandler := solarPanelData.NewDeleteSolarPanelDataHandler(
		solarPanelDataService,
		logger,
	)
	updateSolarPanelDataHandler := solarPanelData.NewUpdateSolarPanelDataHandler(
		solarPanelDataService,
		logger,
	)

	s.router.HandleFunc(
		"/solar-panel-data",
		createSolarPanelDataHandler.CreateSolarPanelDataController,
	).Methods(http.MethodPost)
	s.router.HandleFunc(
		"/solar-panel-data/{id}",
		getSolarPanelDataHandler.GetSolarPanelDataController,
	).Methods(http.MethodGet)
	s.router.HandleFunc(
		"/solar-panel-data/{id}",
		deleteSolarPanelDataHandler.DeleteSolarPanelDataController,
	).Methods(http.MethodDelete)
	s.router.HandleFunc(
		"/solar-panel-data/{id}",
		updateSolarPanelDataHandler.UpdateSolarPanelDataController,
	).Methods(http.MethodPut)
}
