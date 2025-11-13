package teamrepository

const (
	createTeamQuery = `
		INSERT INTO teams (team_name)
			VALUES ($1)
	`

	getTeamWithUsersByNameQuery = `
		SELECT t.team_name, u.user_id, u.username, u.is_active
			FROM teams t
		INNER JOIN users u 
			ON t.team_name = u.team_name
		WHERE t.team_name = $1
	`

	getActiveUserIDFromUserTeamQuery = `
		SELECT (user_id) FROM users 
			WHERE team_name IN (
				SELECT team_name from users WHERE user_id = $1
			) AND user_id != $1 AND is_active = true
	`
)
