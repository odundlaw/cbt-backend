-- name: CreateUser :one
INSERT INTO users (
  full_name,
  email,
  password,
  phone
)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: CreateAdmin :one
INSERT INTO users (
  full_name,
  email,
  password,
  role,
  admin_code,
  phone
)
VALUES ($1, $2, $3, 'ADMIN', $4, $5)
RETURNING *;


-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1;

-- name: ListUsers :many
SELECT *
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;


-- name: UpdateUserRole :one
UPDATE users
SET role = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: UpdateLastLogin :one
UPDATE users
SET last_login = now()
WHERE id = $1
RETURNING *;


-- name: UpdateUserPassword :one
UPDATE users
SET password = $2,
    updated_at = now()
WHERE id = $1
RETURNING *;


-- name: UpdateAdminFields :one
UPDATE users
SET admin_code = $2,
    phone = $3,
    updated_at = now()
WHERE id = $1
RETURNING *;

