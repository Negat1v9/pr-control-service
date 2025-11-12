package teamrepository

const (
	createTeamQuery = `
		INSERT INTO teams (team_name)
			VALUES ($1)
		RETURNING *
	`

	getTeamByNameQuery = `
		SELECT team_name 
			FROM teams
		WHERE team_name = $1
	`

	getTeamWithUsersByNameQuery = `
		SELECT t.team_name, u.user_id, u.username, u.is_active
			FROM teams t
		LEFT JOIN users u 
			ON t.team_name = u.team_name
	`

	getActiveUserFromUserTeamWithExceptionQuery = `
		SELECT (user_id) 
			FROM users 
		WHERE team_name IN (
			SELECT team_name FROM users WHERE user_id = $1
		) AND is_active = true AND user_id NOT IN (%s)
	`
)
