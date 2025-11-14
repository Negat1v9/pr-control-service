package app

import (
	"github.com/Negat1v9/pr-review-service/config"
	prservice "github.com/Negat1v9/pr-review-service/internal/pullRequest/service"
	"github.com/Negat1v9/pr-review-service/internal/server"
	"github.com/Negat1v9/pr-review-service/internal/store"
	teamservice "github.com/Negat1v9/pr-review-service/internal/team/service"
	userservice "github.com/Negat1v9/pr-review-service/internal/users/service"
	"github.com/Negat1v9/pr-review-service/pkg/logger"
	"github.com/Negat1v9/pr-review-service/pkg/postgres"
)

type App struct {
	cfg *config.Config
	log *logger.Logger
}

func New(cfg *config.Config, log *logger.Logger) *App {
	return &App{
		cfg: cfg,
		log: log,
	}
}

func (a *App) Run() error {
	db, err := postgres.NewPostgresConn(a.cfg.PostgresConfig.DbHost, a.cfg.PostgresConfig.DbPort, a.cfg.PostgresConfig.DbUser, a.cfg.PostgresConfig.DbPassword, a.cfg.PostgresConfig.DbName)
	if err != nil {
		return err
	}
	a.log.Infof("connect to postgres host: %s, port %d", a.cfg.PostgresConfig.DbHost, a.cfg.PostgresConfig.DbPort)

	storage := store.NewStore(db)

	teamService := teamservice.NewTeamService(storage)
	userService := userservice.NewUserService(storage)
	prService := prservice.NewPRService(storage)

	server := server.New(a.cfg, a.log)

	server.MapHandlers(teamService, userService, prService)
	return server.Run()
}
