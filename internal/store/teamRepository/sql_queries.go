package teamrepository

const (
	createTeamQuery = `
		INSERT INTO teams (team_name)
			VALUES ($1)
	`

	getTeamWithUsersByNameQuery = `
		SELECT t.team_name, u.user_id, u.username, u.is_active
			FROM team_member tm
		INNER JOIN teams t 
			ON tm.team_name = t.team_name
		INNER JOIN users u 
			ON tm.user_id = u.user_id
		WHERE tm.team_name = $1
	`

	getActiveUserIDFromUserTeamQuery = `
		SELECT (user_id) FROM team_member 
			WHERE team_name IN (
				SELECT team_name from team_member WHERE user_id = $1
			) AND user_id != $1
	`

	createTeamMemberQuery = `
		INSERT INTO team_member (user_id, team_name) 
			VALUES ($1, $2)
	`

	createManyTeamMembersQuery = `
		INSERT INTO team_member (user_id, team_name) 
			VALUES %s
	`
)
