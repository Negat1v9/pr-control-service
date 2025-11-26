package server

import (
	"context"
	"net/http"

	"github.com/Negat1v9/pr-review-service/config"
	"github.com/Negat1v9/pr-review-service/pkg/logger"
	"github.com/Negat1v9/pr-review-service/pkg/metrics"
)

type Server struct {
	log     *logger.Logger
	server  http.Server
	cfg     *config.Config
	metrics *metrics.PrometheusMetrics
}

func New(cfg *config.Config, log *logger.Logger, metricsManager *metrics.PrometheusMetrics) *Server {
	return &Server{
		log: log,
		server: http.Server{
			Addr: cfg.WebConfig.ListenAddress,
		},
		cfg:     cfg,
		metrics: metricsManager,
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
