-- name: FindAllUsers :many
SELECT * FROM users;

-- name: FindUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: CreateUser :exec
INSERT INTO users (id, email, name, avatar_url)
VALUES ($1, $2, $3, $4);

-- name: DeleteUserById :exec
DELETE FROM users
where id = $1;

-- name: UpdateUserEmail :exec
UPDATE users
SET email = $2
WHERE id = $1;

