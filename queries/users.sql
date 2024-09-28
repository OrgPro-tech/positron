-- name: GetUser :one
SELECT email, password from users where email = $1;

-- name: GetUsers :many
SELECT id, name from users where 1;
-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY name;

-- name: CreateUser :one
INSERT INTO users (
  username, password, email, name, mobile_number, user_type, business_id, outlet_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET username = $2, email = $3, name = $4, mobile_number = $5, user_type = $6, 
    business_id = COALESCE($7, business_id), outlet_id = COALESCE($8, outlet_id)
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: GetUserOutletByID :one
SELECT * FROM user_outlets
WHERE id = $1 LIMIT 1;

-- name: ListUserOutlets :many
SELECT * FROM user_outlets
WHERE user_id = $1;

-- name: CreateUserOutlet :one
INSERT INTO user_outlets (
  user_id, business_id, outlet_id
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: DeleteUserOutlet :exec
DELETE FROM user_outlets
WHERE id = $1;

-- name: GetBusinessByID :one
SELECT * FROM businesses
WHERE id = $1 LIMIT 1;

-- name: ListBusinesses :many
SELECT * FROM businesses
ORDER BY company_name;

-- name: CreateBusiness :one
INSERT INTO businesses (
  contact_person_name, contact_person_email, contact_person_mobile_number,
  company_name, address, pin, city, state, country, business_type,
  gst, pan, bank_account_number, bank_name, ifsc_code, account_type, account_holder_name
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
)
RETURNING *;

-- name: UpdateBusiness :one
UPDATE businesses
SET contact_person_name = $2, contact_person_email = $3, contact_person_mobile_number = $4,
    company_name = $5, address = $6, pin = $7, city = $8, state = $9, country = $10,
    business_type = $11, gst = $12, pan = $13, bank_account_number = $14, bank_name = $15,
    ifsc_code = $16, account_type = $17, account_holder_name = $18
WHERE id = $1
RETURNING *;

-- name: DeleteBusiness :exec
DELETE FROM businesses
WHERE id = $1;

-- name: GetOutletByID :one
SELECT * FROM outlets
WHERE id = $1 LIMIT 1;

-- name: ListOutlets :many
SELECT * FROM outlets
WHERE business_id = $1;

-- name: CreateOutlet :one
INSERT INTO outlets (
  outlet_name, outlet_address, outlet_pin, outlet_city, outlet_state,
  outlet_country, business_id
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: UpdateOutlet :one
UPDATE outlets
SET outlet_name = $2, outlet_address = $3, outlet_pin = $4, outlet_city = $5,
    outlet_state = $6, outlet_country = $7, business_id = $8
WHERE id = $1
RETURNING *;

-- name: DeleteOutlet :exec
DELETE FROM outlets
WHERE id = $1;

-- name: GetUserSessionByUserID :one
SELECT * FROM user_sessions
WHERE user_id = $1 LIMIT 1;

-- name: CreateUserSession :one
INSERT INTO user_sessions (
  user_id, access_token, refresh_token, expire_at
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;

-- name: UpdateUserSession :one
UPDATE user_sessions
SET access_token = $2, refresh_token = $3, expire_at = $4
WHERE user_id = $1
RETURNING *;

-- name: DeleteUserSession :exec
DELETE FROM user_sessions
WHERE user_id = $1;