-- name: GetUserById :one
SELECT 
    username,
    email,
    avatar,
    skin,
    cloak,
    registered_at,
    is_active
FROM users
WHERE id = $1;

-- name: CreateUser :many
INSERT INTO users (
    username, 
    email, 
    password_hash, 
    avatar, 
    skin, 
    cloak, 
    registered_at, 
    is_active
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING id;