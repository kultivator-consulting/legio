-- name: GetAccount :one
SELECT
    id,
    created,
    modified,
    deleted,
    is_system,
    first_name,
    last_name,
    access_Level,
    credit_balance,
    username,
    user_email_address,
    user_avatar_url,
    user_zone_info,
    user_locale,
    is_locked,
    last_login,
    reset_password_token,
    reset_password_token_expiry FROM account
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetAccountById :one
SELECT * FROM account
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetUnlockedAccountById :one
SELECT * FROM account
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_locked = FALSE
LIMIT 1;

-- name: GetAccountByUsername :one
SELECT * FROM account
WHERE username = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_locked = FALSE
LIMIT 1;

-- name: GetAccountByEmailAddress :one
SELECT * FROM account
WHERE user_email_address = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_locked = FALSE
LIMIT 1;

-- name: CountAccounts :one
SELECT
    COUNT(id)
FROM account
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListAccounts :many
SELECT
    id,
    created,
    modified,
    deleted,
    is_system,
    first_name,
    last_name,
    access_Level,
    credit_balance,
    username,
    user_email_address,
    user_avatar_url,
    user_zone_info,
    user_locale,
    is_locked,
    last_login,
    reset_password_token,
    reset_password_token_expiry FROM account
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY
    CASE
        WHEN sqlc.arg(sort_by)::text = 'created' AND sqlc.arg(sort_order) = 'asc' THEN sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'created' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'first_name' AND sqlc.arg(sort_order) = 'asc' THEN sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'first_name' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'last_name' AND sqlc.arg(sort_order) = 'asc' THEN sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'last_name' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'access_Level' AND sqlc.arg(sort_order) = 'asc' THEN sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'access_Level' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'credit_balance' AND sqlc.arg(sort_order) = 'asc' THEN sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'credit_balance' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'username' AND sqlc.arg(sort_order) = 'asc' THEN sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'username' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'user_email_address' AND sqlc.arg(sort_order) = 'asc' THEN sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'user_email_address' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'is_locked' AND sqlc.arg(sort_order) = 'asc' THEN  sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'is_locked' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'last_login' AND sqlc.arg(sort_order) = 'asc' THEN sqlc.arg(sort_by)::text END ASC,
    CASE
        WHEN sqlc.arg(sort_by)::text = 'last_login' AND sqlc.arg(sort_order) = 'desc' THEN sqlc.arg(sort_by)::text END DESC

OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateAccount :one
INSERT INTO account (
    id,
    created,
    modified,
    deleted,
    first_name,
    last_name,
    access_Level,
    credit_balance,
    username,
    user_email_address,
    user_avatar_url,
    user_zone_info,
    user_locale,
    user_password,
    is_locked,
    last_login,
    reset_password_token,
    reset_password_token_expiry
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
         )
RETURNING *;

-- name: CreateAccountAndReturnId :one
INSERT INTO account (
    id,
    created,
    modified,
    deleted,
    first_name,
    last_name,
    access_Level,
    credit_balance,
    username,
    user_email_address,
    user_avatar_url,
    user_zone_info,
    user_locale,
    user_password,
    is_locked,
    last_login,
    reset_password_token,
    reset_password_token_expiry
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
         )
RETURNING id;

-- name: DeleteAccount :exec
UPDATE account
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
  AND is_system = false;

-- name: UpdateAccountLastLogin :exec
UPDATE account
SET
    last_login = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateAccount :one
UPDATE account
SET
    modified = now() AT TIME ZONE 'UTC',
    first_name = $2,
    last_name = $3,
    access_Level = $4,
    credit_balance = $5,
    username = $6,
    user_email_address = $7,
    user_avatar_url = $8,
    user_zone_info = $9,
    user_locale = $10,
    user_password = $11,
    is_locked = $12,
    last_login = $13,
    reset_password_token = $14,
    reset_password_token_expiry = $15
WHERE id = $1
RETURNING
    id,
    created,
    modified,
    deleted,
    first_name,
    last_name,
    access_Level,
    credit_balance,
    username,
    user_email_address,
    user_avatar_url,
    user_zone_info,
    user_locale,
    is_locked,
    last_login,
    reset_password_token,
    reset_password_token_expiry;

-- name: UpdateSystemProfile :one
UPDATE account
SET
    modified = now() AT TIME ZONE 'UTC',
    user_password = $2
WHERE id = $1
RETURNING
    id,
    created,
    modified,
    deleted,
    first_name,
    last_name,
    access_Level,
    credit_balance,
    username,
    user_email_address,
    user_avatar_url,
    user_zone_info,
    user_locale,
    is_locked,
    last_login,
    reset_password_token,
    reset_password_token_expiry;

-- name: UpdateAccountResetPasswordToken :one
UPDATE account
SET
    modified = now() AT TIME ZONE 'UTC',
    reset_password_token = $2,
    reset_password_token_expiry = $3
WHERE id = $1
RETURNING
    id,
    created,
    modified,
    deleted,
    first_name,
    last_name,
    access_Level,
    credit_balance,
    username,
    user_email_address,
    user_avatar_url,
    user_zone_info,
    user_locale,
    is_locked,
    last_login,
    reset_password_token,
    reset_password_token_expiry;

-- name: UpdateAccountResetPasswordByToken :one
UPDATE account
SET
    modified = now() AT TIME ZONE 'UTC',
    user_password = $2,
    reset_password_token = NULL
WHERE reset_password_token = $1
    AND reset_password_token_expiry > now() AT TIME ZONE 'UTC'
RETURNING
    id,
    created,
    modified,
    deleted,
    first_name,
    last_name,
    access_Level,
    credit_balance,
    username,
    user_email_address,
    user_avatar_url,
    user_zone_info,
    user_locale,
    is_locked,
    last_login,
    reset_password_token,
    reset_password_token_expiry;
