package userrepository

const (
	createUserQuery = `
		INSERT INTO users (user_id, username, team_name, is_active)
			VALUES ($1, $2, $3, $4)
		RETURNING *
	`
	createManyUsersQuery = `
		INSERT INTO users (user_id, username, team_name, is_active)
			VALUES (%s)
		RETURNING *
	`
	getUserReviewsQuery = `
		SELECT u.user_id, pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at 
			FROM assigned_reviewers ar 
		JOIN users u 
			ON ar.reviewer_user_id = u.user_id 
		JOIN pull_requests pr 
			ON ar.pull_request_id = pr.pull_request_id 
		WHERE ar.reviewer_user_id = $1`
)
