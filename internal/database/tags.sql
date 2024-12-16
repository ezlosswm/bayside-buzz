-- name: CreateTag :one
INSERT INTO tags (name)
VALUES ($1)
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: FindTagByName :one
SELECT *
FROM tags
WHERE name = $1;

-- name: LinkEventToTag :exec
INSERT INTO event_tags (eventId, tagId)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: GetTagsForEvent :many
SELECT tags.id, tags.name
FROM tags
JOIN event_tags ON tags.id = event_tags.tagId
WHERE event_tags.eventId = $1;
