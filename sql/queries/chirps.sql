-- name: AddChirp :one
INSERT INTO chirps (body, user_id)
    VALUES ($1, $2)
RETURNING
    *;

-- name: GetAllChirpsASC :many
SELECT
    *
FROM
    chirps
WHERE (sqlc.narg ('author_id')::UUID IS NULL
    OR user_id = sqlc.narg ('author_id')::UUID)
ORDER BY
    created_at ASC;

-- name: GetAllChirpsDESC :many
SELECT
    *
FROM
    chirps
WHERE (sqlc.narg ('author_id')::UUID IS NULL
    OR user_id = sqlc.narg ('author_id')::UUID)
ORDER BY
    created_at DESC;

-- name: GetChirp :one
SELECT
    *
FROM
    chirps
WHERE
    id = $1;

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE id = $1;

