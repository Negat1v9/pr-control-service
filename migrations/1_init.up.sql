CREATE TABLE IF NOT EXISTS teams (
    team_name TEXT PRIMARY KEY
);


CREATE TABLE IF NOT EXISTS users (
    user_id TEXT PRIMARY KEY,
    username TEXT NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

CREATE INDEX IF NOT EXISTS idx_users_user_id ON users(user_id);

CREATE TABLE IF NOT EXISTS team_member (
    user_id TEXT NOT NULL REFERENCES users(user_id),
    team_name TEXT NOT NULL REFERENCES teams(team_name),
    UNIQUE(user_id, team_name)
);

CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id TEXT REFERENCES users(user_id),
    status VARCHAR(15) DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'MERGED')),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    merged_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_pull_requests_pull_request_id ON pull_requests(pull_request_id);
CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id);


CREATE TABLE IF NOT EXISTS assigned_reviewers (
    id SERIAL PRIMARY KEY,
    reviewer_user_id TEXT NOT NULL REFERENCES users(user_id),
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id),
    UNIQUE(reviewer_user_id, pull_request_id)
);