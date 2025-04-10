-- name: CreateUser :exec
INSERT INTO users (id, pass_hash, username)
VALUES ($1, $2, $3);

-- name: GetUser :one
SELECT pass_hash, username
FROM users
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
