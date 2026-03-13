-- name: GetUserByID :one
SELECT 
    id,
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

-- name: GetUserByEmail :one
SELECT
    id,
    username,
    email,
    avatar,
    skin,
    cloak,
    registered_at,
    is_active
FROM users
WHERE email = $1
LIMIT 1;


-- name: GetUserByUsername :one
SELECT
    id,
    username,
    email,
    avatar,
    skin,
    cloak,
    registered_at,
    is_active
FROM users
WHERE username = $1
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

-- name: GetUserSkin :one
SELECT skin from users WHERE id = $1;

-- name: UpdateUserStatus :one
UPDATE users
SET is_active = $1
WHERE id = $2
RETURNING id;

-- name: UpdateUserSkin :exec
UPDATE users 
SET skin = $1 
WHERE id = $2;

-- name: DeleteUserSkin :exec
UPDATE users
SET skin = NULL
WHERE id = $1;

-- name: GetUserCloak :one
SELECT cloak from users WHERE id = $1;

-- name: UpdateUserCloak :exec
UPDATE users 
SET cloak = $1 
WHERE id = $2;

-- name: DeleteUserCloak :exec
UPDATE users
SET cloak = NULL
WHERE id = $1;

-- name: UpdateUserAvatar :exec
UPDATE users
SET avatar = $1
WHERE id = $2;
