-- name: CreateUser :one
    INSERT INTO "users" ("username","hashed_password","full_name","email") VALUES ($1,$2,$3,$4) RETURNING *;

-- name: GetUser :one
    SELECT * FROM "users" WHERE username=$1 ORDER BY username LIMIT 1;

-- name: GetAllUsers :many
    SELECT * FROM "users" ORDER BY username LIMIT $1 OFFSET $2;

-- name: UpdateUser :one
    UPDATE users
    SET
        hashed_password = coalesce(sqlc.narg(hashed_password), hashed_password),
        password_changed_at = coalesce(sqlc.narg(password_changed_at), password_changed_at),
        full_name = coalesce(sqlc.narg(full_name), full_name),
        email = coalesce(sqlc.narg(email), email),
        is_email_verified = coalesce(sqlc.narg(is_email_verified), is_email_verified)
    WHERE
        username = sqlc.arg(username)
    RETURNING *;

-- name: DeleteUser :exec
    DELETE FROM "users" WHERE username=$1;