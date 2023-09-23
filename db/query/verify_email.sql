-- name: UpdateVerifyEmail :one
UPDATE "verify_emails" SET "username"=$1 WHERE "username"=$2 RETURNING *;

-- name: CreateVerifyEmail :one
INSERT INTO "verify_emails" (username,email,secret_code) VALUES($1,$2,$3) RETURNING *;

-- name: GetVerifyEmail :one
SELECT * FROM "verify_emails" WHERE "username"=$1;