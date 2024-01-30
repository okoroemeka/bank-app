-- name: UpdateVerifyEmail :one
UPDATE "verify_emails"
    SET
        "username" = coalesce(sqlc.narg(username), "username"),
        "email" = coalesce(sqlc.narg(email), "email"),
        "secret_code" = coalesce(sqlc.narg(secret_code), "secret_code"),
        "is_used" = coalesce(sqlc.narg(is_used), "is_used")
    WHERE
        "id"=sqlc.arg(id)
    RETURNING *;

-- name: UpdateVerifyEmailIsUsedField :one
UPDATE "verify_emails"
    SET
        "is_used"=$1
    WHERE
        "id"=$2
    AND
        "secret_code"=$3
    AND
    "expired_at">now()
    AND NOT
        "is_used"
    RETURNING *;

-- name: CreateVerifyEmail :one
INSERT INTO "verify_emails" (username,email,secret_code) VALUES($1,$2,$3) RETURNING *;

-- name: GetVerifyEmail :one
SELECT * FROM "verify_emails" WHERE "id"=$1;