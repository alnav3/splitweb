-- name: FindUserById :one
SELECT *
FROM users
WHERE id = $1;
-- name: CreateUser :exec
INSERT INTO users (id, email, name, avatar_url)
VALUES ($1, $2, $3, $4);

