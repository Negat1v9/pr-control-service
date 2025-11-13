package server

import (
	"net/http"

	teamhttp "github.com/Negat1v9/pr-review-service/internal/team/http"
	teamservice "github.com/Negat1v9/pr-review-service/internal/team/service"
	userhttp "github.com/Negat1v9/pr-review-service/internal/users/http"
	userservice "github.com/Negat1v9/pr-review-service/internal/users/service"
)

func (s *Server) MapHandlers(teamService *teamservice.TeamService, userService *userservice.UserService) {
	router := http.NewServeMux()

	teamHandler := teamhttp.NewTeamHanlder(s.log, teamService)
	userHandler := userhttp.NewUserHandler(s.log, userService)

	teamRouter := teamhttp.TeamRouter(teamHandler)
	userRouter := userhttp.UserRouter(userHandler)

	router.Handle("/team/", http.StripPrefix("/team", teamRouter))
	router.Handle("/users/", http.StripPrefix("/users", userRouter))

	s.server.Handler = router
}
