-- name: CreateUser :one
INSERT INTO users (
  username,
  hashed_password
) VALUES (
  $1, $2
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE username = $1 
LIMIT 1;

-- name: UpdateUser :one
UPDATE users
SET
  hashed_password = @hashed_password
WHERE
  username = @username
RETURNING *;
