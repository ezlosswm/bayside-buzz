-- name: CreateUser :exec
INSERT INTO users (
    name, email, password_hash
    ) VALUES (
    $1, $2, $3
);

-- name: GetUsers :many 
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE email = $1;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: FindUser :one
SELECT * FROM users WHERE id = $1;
