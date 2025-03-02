-- name: CreateUser :one
INSERT INTO users (id, email)
VALUES (
    $1,
    $2
)
RETURNING id;

-- name: SetUserEmailVerified :one
UPDATE users
SET email_verified = $2
WHERE id = $1
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;
