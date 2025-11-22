CREATE TABLE teams (
    team_name TEXT PRIMARY KEY
);

CREATE TABLE users (
    user_id     TEXT PRIMARY KEY,
    username    TEXT NOT NULL,
    team_name   TEXT NOT NULL REFERENCES teams(team_name),
    is_active   BOOLEAN NOT NULL
);

CREATE TABLE pull_requests (
    pull_request_id   TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id         TEXT NOT NULL REFERENCES users(user_id),
    status            TEXT NOT NULL,
    assigned_reviewers TEXT[] NOT NULL,
    created_at        TIMESTAMPTZ NOT NULL,
    merged_at         TIMESTAMPTZ 
);

CREATE INDEX idx_users_team ON users(team_name);
CREATE INDEX idx_pull_requests_reviewers ON pull_requests USING GIN (assigned_reviewers);