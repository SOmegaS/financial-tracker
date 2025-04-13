-- name: GetReport :many
SELECT category, amount
FROM bills
WHERE user_id = $1;
-- AND tmstmp >= $2
-- AND tmstmp < $3;
