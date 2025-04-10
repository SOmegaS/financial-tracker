-- name: CreateCategory :exec
INSERT INTO categories (user_id, name)
VALUES ($1, $2);

-- name: ListCategories :many
SELECT name
FROM categories
WHERE user_id = $1;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE user_id = $1
AND name = $2;

-- name: DeleteCategoriesByUser :exec
DELETE FROM categories
WHERE user_id = $1;
