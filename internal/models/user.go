package models

type User struct {
	UserID   string `json:"user_id" db:"user_id"`
	Username string `json:"username" db:"username"`
	TeamName string `json:"team_name,omitempty" db:"team_name"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type SetUserActiveStatusRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
