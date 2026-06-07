-- name: EmptyUsers :exec
TRUNCATE TABLE users RESTART IDENTITY CASCADE;

-- name: EmptyChirps :exec
TRUNCATE TABLE chirps RESTART IDENTITY CASCADE;