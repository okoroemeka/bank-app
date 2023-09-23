-- name: CreateUser :one
    INSERT INTO "users" ("username","hashed_password","full_name","email") VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetUser :one
    SELECT * FROM "users" WHERE username=$1 ORDER BY username LIMIT 1;

-- name: GetAllUsers :many
    SELECT * FROM "users" ORDER BY username LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
    UPDATE "users" SET "full_name"=$1, "email"=$2 WHERE username=$3 RETURNING *;

-- name: UpdateUserPassword :one
UPDATE "users" SET "hashed_password"=$1, "password_changed_at"=(now()) WHERE username=$2 RETURNING *;

-- name: DeleteUser :exec
    DELETE FROM "users" WHERE username=$1;