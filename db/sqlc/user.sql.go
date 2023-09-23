// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.21.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
    INSERT INTO "users" ("username","hashed_password","full_name","email") VALUES ($1,$2,$3,$4) RETURNING username, hashed_password, full_name, email, is_email_verified, password_changed_at, created_at
`

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.Username,
		arg.HashedPassword,
		arg.FullName,
		arg.Email,
	)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.IsEmailVerified,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
    DELETE FROM "users" WHERE username=$1
`

func (q *Queries) DeleteUser(ctx context.Context, username string) error {
	_, err := q.db.ExecContext(ctx, deleteUser, username)
	return err
}

const getAllUsers = `-- name: GetAllUsers :many
    SELECT username, hashed_password, full_name, email, is_email_verified, password_changed_at, created_at FROM "users" ORDER BY username LIMIT $1 OFFSET $2
`

type GetAllUsersParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) GetAllUsers(ctx context.Context, arg GetAllUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getAllUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.Username,
			&i.HashedPassword,
			&i.FullName,
			&i.Email,
			&i.IsEmailVerified,
			&i.PasswordChangedAt,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
    SELECT username, hashed_password, full_name, email, is_email_verified, password_changed_at, created_at FROM "users" WHERE username=$1 ORDER BY username LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.IsEmailVerified,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
    UPDATE "users" SET "full_name"=$1, "email"=$2 WHERE username=$3 RETURNING username, hashed_password, full_name, email, is_email_verified, password_changed_at, created_at
`

type UpdateUserParams struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser, arg.FullName, arg.Email, arg.Username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.IsEmailVerified,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :one
UPDATE "users" SET "hashed_password"=$1, "password_changed_at"=(now()) WHERE username=$2 RETURNING username, hashed_password, full_name, email, is_email_verified, password_changed_at, created_at
`

type UpdateUserPasswordParams struct {
	HashedPassword string `json:"hashed_password"`
	Username       string `json:"username"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPassword, arg.HashedPassword, arg.Username)
	var i User
	err := row.Scan(
		&i.Username,
		&i.HashedPassword,
		&i.FullName,
		&i.Email,
		&i.IsEmailVerified,
		&i.PasswordChangedAt,
		&i.CreatedAt,
	)
	return i, err
}
