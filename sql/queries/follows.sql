-- name: CreateFollow :one
INSERT INTO follows (id, feed_id, user_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: DeleteFollow :exec
DELETE FROM follows
    WHERE id=$1;
