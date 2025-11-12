package pullrequestrepository

const (
	createPullRequestQuery = `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id)
			VALUES ($1, $2, $3)
		RETURNING *
	`

	getPullRequestWithAssignedsByIDQuery = `
		SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, pr.status, pr.created_at, pr.merged_at, ar.reviewer_user_id 
			FROM pull_requests pr 
		LEFT JOIN assigned_reviewers ar 
			ON ar.pull_request_id = pr.pull_request_id
		WHERE pr.pull_request_id = $1

	`

	updatePullRequestQuery = `
		WITH updated AS (
			UPDATE pull_requests
				SET status = COALESCE(NULLIF($1, ''), status),
				merged_at = COALESCE(NULLIF($2, null), merged_at)
			WHERE pull_request_id = $3
			RETURNING pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		)
	SELECT u.pull_request_id, u.pull_request_name, u.author_id, u.status, u.created_at, u.merged_at, ar.reviewer_user_id 
		FROM updated u 
	JOIN assigned_reviewers ar ON ar.pull_request_id = u.pull_request_id
	`

	createAssignedQuery = `
		INSERT INTO assigned_reviewers (reviewer_user_id, pull_request_id) VALUES ($1, $2)
	`
)
