-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, user_id, body)
VALUES (gen_random_uuid(), NOW(), NOW(), $1, $2)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY created_at;

-- name: GetChirpsByUserId :many
SELECT * FROM chirps WHERE user_id = $1 ORDER BY created_at;

-- name: GetChirpById :one
SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirpById :exec
DELETE FROM chirps WHERE id = $1;

-- name: ResetChirps :exec
DELETE FROM chirps;