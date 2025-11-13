package server

import (
	"net/http"

	teamhttp "github.com/Negat1v9/pr-review-service/internal/team/http"
	teamservice "github.com/Negat1v9/pr-review-service/internal/team/service"
)

func (s *Server) MapHandlers(teamService *teamservice.TeamService) {
	router := http.NewServeMux()

	teamHandler := teamhttp.NewTeamHanlder(s.log, teamService)

	teamRouter := teamhttp.TeamRouter(teamHandler)

	router.Handle("/team/", http.StripPrefix("/team", teamRouter))

	s.server.Handler = router
}
