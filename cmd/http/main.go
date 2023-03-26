package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/loukaspe/solar-panel-data-crud/internal/repositories"
	"github.com/loukaspe/solar-panel-data-crud/pkg/server"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	logger := log.New()

	err := godotenv.Load("./config/.env")
	if err != nil {
		logger.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Fatal("Error starting service")
	}

	router := mux.NewRouter()
	db := make(repositories.SolarPanelDataDB)
	httpServer := &http.Server{
		Addr:    os.Getenv("SERVER_ADDR"),
		Handler: router,
	}

	server := server.NewServer(db, router, httpServer, logger)

	server.Run()
}
