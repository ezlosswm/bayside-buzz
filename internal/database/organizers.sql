-- name: CreateOrganizer :exec
INSERT INTO organizers (
    organizer_name, description, img_url
    ) VALUES (
    $1, $2, $3
);

-- name: GetOrganizers :many 
SELECT * FROM organizers;

-- name: GetOrganizer :one
SELECT * FROM organizers WHERE id = $1;

-- name: CountOrganizers :one
SELECT COUNT(*) FROM organizers;

-- name: DeleteOrganizer :exec
DELETE FROM organizers
WHERE id = $1;

-- name: UpdateOrganizer :exec
UPDATE organizers 
SET 
    organizer_name = $1, 
    description = $2,
    img_url = $3
WHERE 
    id = $4;