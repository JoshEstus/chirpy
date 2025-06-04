-- name: CreateChirp :one
INSERT INTO chirps(id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
SELECT *
FROM chirps
where 1=1
order by created_at ASC;

-- name: GetChirp :one
SELECT *
FROM chirps
WHERE 1 = 1
AND id = $1;