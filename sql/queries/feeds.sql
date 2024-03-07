-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT id, created_at, updated_at, name, url, user_id, last_fetched_at
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = $2, updated_at = $2
WHERE id = $1
RETURNING *;
