-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2,
    FALSE
)
RETURNING *;

-- name: GetUserById :one
SELECT *
FROM users 
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT *
FROM users 
WHERE email = $1;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: UpdateUser :one
UPDATE users
SET email = $2, hashed_password = $3
WHERE id = $1

RETURNING *;

-- name: UpdateUserChirpyRed :one
UPDATE users
SET is_chirpy_red = $2
WHERE id = $1

RETURNING *;