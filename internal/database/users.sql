-- name: CreateUser :exec
INSERT INTO users (
    name, email, password_hash
) VALUES (
  ?, ?, ?
);

-- name: GetUsers :many 
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users WHERE email = ?;

-- name: CountUsers :one
SELECT COUNT(*) FROM users;

-- name: FindUser :one
SELECT * FROM users WHERE id = ?;
