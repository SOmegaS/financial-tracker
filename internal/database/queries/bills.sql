-- name: CreateBill :exec
INSERT INTO bills (id, user_id, amount, category_name, ts)
VALUES (gen_random_uuid(), $1, $2, $3, $4);

-- name: ListBills :many
SELECT amount, category_name
FROM bills
WHERE user_id = $1
AND ts > $2
AND ts < $3;

-- name: DeleteBill :exec
DELETE FROM bills
WHERE id = $1;

-- name: DeleteBillsByUser :exec
DELETE FROM bills
WHERE user_id = $1;

-- name: DeleteBillsByCategory :exec
DELETE FROM bills
WHERE user_id = $1
AND category_name = $2;
