-- name: CreateUser :one
INSERT INTO users (email, hashed_password)
    VALUES ($1, $2)
RETURNING
    *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT
    *
FROM
    users
WHERE
    email = $1;

-- name: UpdateUserEmailPassword :one
UPDATE
    users
SET
    email = $2,
    hashed_password = $3
WHERE
    id = $1
RETURNING
    *;

-- name: UpgradeToRed :one
UPDATE
    users
SET
    is_chirpy_red = TRUE
WHERE
    id = $1
RETURNING
    *;

