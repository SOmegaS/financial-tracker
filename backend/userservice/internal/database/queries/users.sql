-- name: CreateUser :exec
INSERT INTO users (id, username, password)
VALUES ($1, $2, $3);

-- name: GetUserIdPassword :one
SELECT id, password
FROM users
WHERE username = $1;
