-- name: GetUser :one
SELECT email, password from users where email = $1;

-- name: GetUsers :many
SELECT id, name from users where 1;