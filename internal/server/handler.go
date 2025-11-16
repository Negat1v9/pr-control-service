package server

import (
	"net/http"

	"github.com/Negat1v9/pr-review-service/internal/middleware"
	prhttp "github.com/Negat1v9/pr-review-service/internal/pullRequest/http"
	prservice "github.com/Negat1v9/pr-review-service/internal/pullRequest/service"
	teamhttp "github.com/Negat1v9/pr-review-service/internal/team/http"
	teamservice "github.com/Negat1v9/pr-review-service/internal/team/service"
	userhttp "github.com/Negat1v9/pr-review-service/internal/users/http"
	userservice "github.com/Negat1v9/pr-review-service/internal/users/service"
)

func (s *Server) MapHandlers(teamService *teamservice.TeamService, userService *userservice.UserService, prService *prservice.PRService) {
	router := http.NewServeMux()

	teamHandler := teamhttp.NewTeamHanlder(s.log, teamService)
	userHandler := userhttp.NewUserHandler(s.log, userService)
	prHandler := prhttp.NewPRHanlder(s.log, prService)

	teamRouter := teamhttp.TeamRouter(teamHandler)
	userRouter := userhttp.UserRouter(userHandler)
	prRouter := prhttp.PRRouter(prHandler)

	router.Handle("/team/", http.StripPrefix("/team", teamRouter))
	router.Handle("/users/", http.StripPrefix("/users", userRouter))
	router.Handle("/pullRequest/", http.StripPrefix("/pullRequest", prRouter))

	// middleware service
	mw := middleware.New()

	// all requests go through from basic middleware
	s.server.Handler = mw.BasicMW()(router)
}
