// Code generated by sqlc. DO NOT EDIT.
// source: user.sql

package internal_sql

import (
	"context"
)

const getHashedPwdAndIdfromUsername = `-- name: GetHashedPwdAndIdfromUsername :one
SELECT app_user.id, app_user.hashedPwd from app_user
	WHERE app_user.username = $1
`

type GetHashedPwdAndIdfromUsernameRow struct {
	ID        int64
	Hashedpwd string
}

func (q *Queries) GetHashedPwdAndIdfromUsername(ctx context.Context, username string) (GetHashedPwdAndIdfromUsernameRow, error) {
	row := q.db.QueryRowContext(ctx, getHashedPwdAndIdfromUsername, username)
	var i GetHashedPwdAndIdfromUsernameRow
	err := row.Scan(&i.ID, &i.Hashedpwd)
	return i, err
}

const getIDFromEmail = `-- name: GetIDFromEmail :one
SELECT app_user.id FROM app_user
	WHERE app_user.email = $1
limit 1
`

func (q *Queries) GetIDFromEmail(ctx context.Context, email string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getIDFromEmail, email)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getIDfromUsername = `-- name: GetIDfromUsername :one
SELECT app_user.id FROM app_user 
	WHERE app_user.username = $1 
limit 1
`

func (q *Queries) GetIDfromUsername(ctx context.Context, username string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getIDfromUsername, username)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO app_user (email, username, hashedPwd) values ($1, $2, $3) RETURNING app_user.id
`

type InsertUserParams struct {
	Email     string
	Username  string
	Hashedpwd string
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, insertUser, arg.Email, arg.Username, arg.Hashedpwd)
	var id int64
	err := row.Scan(&id)
	return id, err
}