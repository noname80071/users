-- name: GetUserByID :one
SELECT 
    username,
    email,
    avatar,
    skin,
    cloak,
    registered_at,
    is_active
FROM users
WHERE id = $1
LIMIT 1;

-- name: CreateUser :one
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

-- name: UpdateUserSkin :exec
UPDATE users 
SET skin = $1 
WHERE id = $2;

-- name: UpdateUserCloak :exec
UPDATE users 
SET cloak = $1 
WHERE id = $2;