-- name: CreateFeedFollow :many

WITH cte_feedFollows AS (
    INSERT INTO feed_follows(id,created_at,updated_at,user_id,feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    ) RETURNING *
) SELECT
    cte_feedFollows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM cte_feedFollows
INNER JOIN users ON users.id = cte_feedFollows.user_id
INNER JOIN feeds ON feeds.id = cte_feedFollows.feed_id;

-- name: GetFeedFollowsForUser :many

SELECT 
    feed_follows.*,
    users.name AS user_name,
    feeds.name AS feed_name
FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE user_id = $1
AND feed_id = $2;
