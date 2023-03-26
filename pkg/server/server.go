package server

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/loukaspe/solar-panel-data-crud/internal/repositories"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
	db         repositories.SolarPanelDataDB
	logger     *log.Logger
}

func NewServer(
	db repositories.SolarPanelDataDB,
	router *mux.Router,
	httpServer *http.Server,
	logger *log.Logger,
) *Server {
	return &Server{
		router:     router,
		db:         db,
		httpServer: httpServer,
		logger:     logger,
	}
}

func (s *Server) Run() {
	s.initializeRoutes(s.db, s.logger)

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
