package handlers

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

type HealthCheckHandler struct {
	logger *log.Logger
}

func NewHealthCheckHandler(logger *log.Logger) *HealthCheckHandler {
	return &HealthCheckHandler{logger: logger}
}

func (handler *HealthCheckHandler) HealthCheckController(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(`OK`))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		handler.logger.WithFields(log.Fields{
			"errorMessage": err.Error(),
		}).Error("Error in health check")

		return
	}
}
