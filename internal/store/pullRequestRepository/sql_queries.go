package pullrequestrepository

const (
	createPullRequestQuery = `
		INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id)
			VALUES ($1, $2, $3);
	`

	getPullRequestByIDQuery = `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	getPullRequestReviewersQuery = `
		SELECT reviewer_user_id
		FROM assigned_reviewers
		WHERE pull_request_id = $1
		ORDER BY reviewer_user_id
	`

	mergePullRequestQuery = `
		UPDATE pull_requests
				SET status = 'MERGED',
				merged_at = now()
			WHERE pull_request_id = $1
	`

	getPullRequestsQuantityAssignedReviewers = `
		SELECT pr.pull_request_id, COUNT(ar.reviewer_user_id) AS quantity_reviewers
			FROM pull_requests pr 
		LEFT JOIN assigned_reviewers ar ON pr.pull_request_id = ar.pull_request_id 
			GROUP by pr.pull_request_id
	`
	createAssignedQuery = `
		INSERT INTO assigned_reviewers (reviewer_user_id, pull_request_id) VALUES ($1, $2)
		ON CONFLICT (reviewer_user_id, pull_request_id) DO NOTHING
	`
	createManyAssignedQuery = `
		INSERT INTO assigned_reviewers (reviewer_user_id, pull_request_id) VALUES %s
		ON CONFLICT (reviewer_user_id, pull_request_id) DO NOTHING
	`
	deleteAssignedByReviewerIDQuery = `
		DELETE FROM assigned_reviewers WHERE reviewer_user_id = $1
	`
)
