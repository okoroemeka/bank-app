-- name: CreateEntry :one
INSERT INTO "entries" ("account_id","amount") VALUES ($1,$2) RETURNING *;

-- name: UpdateEntry :one
UPDATE "entries" SET "amount"=$1 WHERE id=$2 RETURNING *;

-- name: GetEntry :one
SELECT * FROM "entries" WHERE id=$1;

-- name: GetEntries :many
SELECT * FROM "entries" ORDER BY id LIMIT $1 OFFSET $2;

-- name: DeleteEntry :exec
DELETE FROM "entries" WHERE id=$1;