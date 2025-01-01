-- name: CreateEvent :one
INSERT INTO events (
    title, description, date, freq, organizer, imgPath, userId
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetEvents :many
SELECT * FROM events;

-- name: GetEvent :one
SELECT * FROM events WHERE id = $1;

-- name: GetEventsByOrganizer :many
SELECT * FROM events WHERE organizer = $1;

-- name: GetEventWithTags :one
SELECT
    e.id AS eventId,
    e.title,
    e.description,
    e.date,
    e.freq,
    e.organizer,
    e.imgPath,
    e.userId,
    COALESCE(string_agg(t.name, ','), '') AS tags
FROM
    events e
LEFT JOIN
    event_tags et ON e.id = et.eventId
LEFT JOIN
    tags t ON et.tagId = t.id
WHERE
    e.id = $1
GROUP BY
    e.id, e.title, e.description, e.date, e.freq, e.organizer, e.imgPath, e.userId;

-- name: CountEvents :one
SELECT COUNT(*) FROM events;

-- name: GetEventsWithTags :many
SELECT
    e.id AS eventId,
    e.title,
    e.description,
    e.date,
    e.freq,
    e.organizer,
    e.imgPath,
    e.userId,
    COALESCE(string_agg(t.name, ','), '') AS tags
FROM
    events e
LEFT JOIN
    event_tags et ON e.id = et.eventId
LEFT JOIN
    tags t ON et.tagId = t.id
GROUP BY
    e.id, e.title, e.description, e.date, e.freq, e.organizer, e.imgPath, e.userId;

-- name: UpdateEvent :exec
UPDATE events
SET
    title = $2,
    description = $3,
    date = $4,
    freq = $5,
    imgPath = $6
WHERE id = $1;

-- name: DeleteEvent :exec
DELETE FROM events WHERE id = $1;