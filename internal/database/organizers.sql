-- name: CreateOrganizer :exec
INSERT INTO organizers (
    organizer_name, description, value, img_url
    ) VALUES (
    ?, ?, ?, ?
);

-- name: GetOrganizers :many 
SELECT * FROM organizers;

-- name: CountOrganizers :one
SELECT COUNT(*) FROM organizers;

-- name: DeleteOrganizer :exec
DELETE FROM organizers
WHERE id = ?;
