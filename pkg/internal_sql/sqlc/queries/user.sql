-- name: GetIDfromUsername :one
SELECT app_user.id FROM app_user 
	WHERE app_user.username = $1 
limit 1;

-- name: GetIDFromEmail :one
SELECT app_user.id FROM app_user
	WHERE app_user.email = $1
limit 1;

-- name: GetHashedPwdAndIdfromUsername :one
SELECT app_user.id, app_user.hashedPwd from app_user
	WHERE app_user.username = $1;

-- name: InsertUser :one
INSERT INTO app_user (email, username, hashedPwd) values ($1, $2, $3) RETURNING app_user.id;


