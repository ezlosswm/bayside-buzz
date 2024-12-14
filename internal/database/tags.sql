-- name: CreateTag :one
INSERT INTO tags (name)
VALUES (?)
ON CONFLICT (name) DO NOTHING
RETURNING *
;

-- name: FindTagByName :one
SELECT *
FROM tags
WHERE name = ?;

-- name: LinkEventToTag :exec
INSERT INTO event_tags (eventId, tagId)
VALUES (?, ?)
ON CONFLICT DO NOTHING;

-- name: GetTagsForEvent :many
SELECT tags.id, tags.name
FROM tags
JOIN event_tags ON tags.id = event_tags.tagId
WHERE event_tags.eventId = ?;
