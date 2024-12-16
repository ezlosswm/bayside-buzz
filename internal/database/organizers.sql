-- name: CreateOrganizer :exec
INSERT INTO organizers (
    organizer_name, description, img_url
    ) VALUES (
    $1, $2, $3
);

-- name: GetOrganizers :many 
SELECT * FROM organizers;

-- name: CountOrganizers :one
SELECT COUNT(*) FROM organizers;

-- name: DeleteOrganizer :exec
DELETE FROM organizers
WHERE id = $1;
