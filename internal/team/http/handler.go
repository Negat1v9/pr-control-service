package teamhttp

import "net/http"

type TeamHanler struct {
	// log
}

func NewTeamHanlder() *TeamHanler {
	return &TeamHanler{}

}

func (h *TeamHanler) Add(w http.ResponseWriter, r *http.Request) {

}
func (h *TeamHanler) Get(w http.ResponseWriter, r *http.Request) {

}
