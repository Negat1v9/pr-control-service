package models

type Team struct {
	TeamName string `json:"team_name" db:"team_name"`
	Members  []User `json:"members" db:"members"`
}
