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
RETURNING id;

-- name: GetUserByUsernameOrEmail :one
SELECT id, username, password, email, name, user_type, business_id
FROM users
WHERE username = $1 OR email = $1
LIMIT 1;

-- name: CreateUserWithBusiness :one
WITH inserted_business AS (
    INSERT INTO businesses (
        contact_person_name, contact_person_email, contact_person_mobile_number,
        company_name, address, pin, city, state, country, business_type,
        gst, pan, bank_account_number, bank_name, ifsc_code, account_type, account_holder_name
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
    ) RETURNING id
)
INSERT INTO users (
    username, password, email, name, mobile_number, user_type, business_id
) VALUES (
    $18, $19, $20, $21, $22, $23, 
    (SELECT id FROM inserted_business)
) RETURNING id, username, email, name, mobile_number, user_type, business_id;
   
-- name: CreateUserWithBusinessAndOutlets :one
   WITH inserted_business AS (
    INSERT INTO businesses (
        contact_person_name, contact_person_email, contact_person_mobile_number,
        company_name, address, pin, city, state, country, business_type,
        gst, pan, bank_account_number, bank_name, ifsc_code, account_type, account_holder_name
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
    ) RETURNING id
), inserted_outlet AS (
    INSERT INTO outlets (
        outlet_name, outlet_address, outlet_pin, outlet_city, outlet_state, outlet_country, business_id
    ) VALUES (
        $18, $19, $20, $21, $22, $23, (SELECT id FROM inserted_business)
    ) RETURNING id
)
INSERT INTO users (
    username, password, email, name, mobile_number, user_type, business_id, outlet_id
) VALUES (
    $24, $25, $26, $27, $28, $29::UserType, 
    (SELECT id FROM inserted_business), 
    (SELECT id FROM inserted_outlet)
) RETURNING id, username, email, name, mobile_number, user_type, business_id, outlet_id;




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




-- name: CreateOutletWithUserAssociation :one
WITH new_outlet AS (
    INSERT INTO outlets (
        outlet_name,
        outlet_address,
        outlet_pin,
        outlet_city,
        outlet_state,
        outlet_country,
        business_id
    ) VALUES (
        $1, $2, $3, $4, $5, $6, $7
    ) RETURNING *
), new_user_outlet AS (
    INSERT INTO user_outlets (
        user_id,
        business_id,
        outlet_id
    ) VALUES (
        $8,
        (SELECT business_id FROM new_outlet),
        (SELECT id FROM new_outlet)
    ) RETURNING *
)
SELECT 
    o.id,
    o.outlet_name,
    o.outlet_address,
    o.outlet_pin,
    o.outlet_city,
    o.outlet_state,
    o.outlet_country,
    o.business_id,
    uo.id AS user_outlet_id
FROM new_outlet o
JOIN new_user_outlet uo ON o.id = uo.outlet_id;


-- name: UpdateOutlet :one
UPDATE outlets
SET 
    outlet_name = COALESCE($2, outlet_name),
    outlet_address = COALESCE($3, outlet_address),
    outlet_pin = COALESCE($4, outlet_pin),
    outlet_city = COALESCE($5, outlet_city),
    outlet_state = COALESCE($6, outlet_state),
    outlet_country = COALESCE($7, outlet_country)
WHERE id = $1
RETURNING *;



-- name: GetLatestUserSession :one
SELECT * FROM user_sessions
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteUserSessions :exec
DELETE FROM user_sessions
WHERE user_id = $1;

-- name: GetUserSessionByRefreshToken :one
SELECT * FROM user_sessions WHERE refresh_token = $1 LIMIT 1;


-- name: GetUserProfile :one
SELECT 
    u.id,
    u.name,
    u.email,
    u.mobile_number,
    u.user_type,
    u.username,
    b.id AS business_id,
    b.company_name,
    b.contact_person_name,
    b.contact_person_email,
    b.contact_person_mobile_number,
    b.address,
    b.pin,
    b.city,
    b.state,
    b.country,
    b.business_type,
    b.gst,
    b.pan
FROM 
    users u
JOIN 
    businesses b ON u.business_id = b.id
WHERE 
    u.id = $1;

-- name: GetAllCategories :many
SELECT * FROM categories
WHERE business_id = $1
ORDER BY name;

-- name: CreateCategory :one
INSERT INTO categories (name, description, business_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CreateMenuItem :one
INSERT INTO menu_items (
    category_id, name, description, price, is_vegetarian, spice_level, is_available, business_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;


-- name: GetAllMenuItemsByBusinessID :many
SELECT 
    mi.id,
    mi.category_id,
    mi.name,
    mi.description,
    mi.price,
    mi.is_vegetarian,
    mi.spice_level,
    mi.is_available,
    mi.is_deleted,
    c.name AS category_name
FROM 
    menu_items mi
LEFT JOIN 
    categories c ON mi.category_id = c.id
WHERE 
    mi.business_id = $1 AND mi.is_deleted = false
ORDER BY 
    mi.name;


-- name: CreateCustomer :one
INSERT INTO customer (phone_number, name, whatsapp, email, address, outlet_id, business_id)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetCustomerByID :one
SELECT * FROM customer WHERE id = $1;

-- name: GetCustomersByBusinessID :many
SELECT * FROM customer WHERE business_id = $1;

-- name: GetCustomersByOutletId :many
SELECT * FROM customer WHERE outlet_id = $1;

-- name: UpdateCustomer :one
UPDATE customer
SET phone_number = $2, name = $3, whatsapp = $4, email = $5, address = $6, outlet_id = $7
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :exec
DELETE FROM customer WHERE id = $1;