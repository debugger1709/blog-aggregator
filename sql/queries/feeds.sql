-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetAllFeeds :many
SELECT users.name as user, feeds.name, feeds.url FROM
users JOIN feeds ON users.id = feeds.user_id;

-- name: GetFeed :one
SELECT * FROM feeds WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET
last_fetched_at = $1,
updated_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT feeds.* FROM
feeds JOIN feed_follows ON feeds.id = feed_follows.feed_id 
WHERE feed_follows.user_id = $1
ORDER BY last_fetched_at ASC NULLS FIRST;