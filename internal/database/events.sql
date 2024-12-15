-- name: CreateEvent :one
INSERT INTO events (
    title, description, date, freq, organizer, imgPath, userId
    ) VALUES (
    ?, ?, ?, ?, ?, ?, ?
) RETURNING *;

-- name: GetEvents :many
SELECT * FROM events;

-- name: GetEvent :one
SELECT * FROM events WHERE id = ?;

-- name: GetEventsByOrganizer :many
SELECT * FROM events WHERE organizer = ?;

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
    COALESCE(GROUP_CONCAT(t.name, ','), '') AS tags
FROM 
    events e
LEFT JOIN 
    event_tags et ON e.id = et.eventId
LEFT JOIN 
    tags t ON et.tagId = t.id
WHERE 
    e.id = ?
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
    COALESCE(GROUP_CONCAT(t.name, ','), '') AS tags
FROM
    events e
LEFT JOIN
    event_tags et ON e.id = et.eventId
LEFT JOIN
    tags t ON et.tagId = t.id
GROUP BY
    e.id;
