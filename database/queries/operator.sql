-- name: GetOperatorById :one
SELECT * FROM operator
WHERE operator.id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetOperatorByOperatorName :one
SELECT * FROM operator
WHERE operator.operator_name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: CountOperators :one
SELECT
    COUNT(id)
FROM operator
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListOperatorsAsc :many
SELECT * FROM operator
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListOperatorsDesc :many
SELECT * FROM operator
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListOperatorsByPassCode :many
SELECT * FROM operator
WHERE is_pass_code = TRUE
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: CreateOperator :one
INSERT INTO operator (
    id,
    created,
    modified,
    deleted,
    operator_name,
    location_id,
    operator_code,
    instruction_cust,
    email_address,
    operator_description,
    operator_bio,
    video_url,
    star_rating,
    operator_image,
    website_url,
    zone_info,
    locale,
    fare_url,
    phone,
    is_pass_code,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
         )
RETURNING *;

-- name: CreateOperatorAndReturnId :one
INSERT INTO operator (
    id,
    created,
    modified,
    deleted,
    operator_name,
    location_id,
    operator_code,
    instruction_cust,
    email_address,
    operator_description,
    operator_bio,
    video_url,
    star_rating,
    operator_image,
    website_url,
    zone_info,
    locale,
    fare_url,
    phone,
    is_pass_code,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
         )
RETURNING id;

-- name: DeleteOperator :exec
UPDATE operator
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateOperator :one
UPDATE operator
SET
    modified = now() AT TIME ZONE 'UTC',
    operator_name = $2,
    location_id = $3,
    operator_code = $4,
    instruction_cust = $5,
    email_address = $6,
    operator_description = $7,
    operator_bio = $8,
    video_url = $9,
    star_rating = $10,
    operator_image = $11,
    website_url = $12,
    zone_info = $13,
    locale = $14,
    fare_url = $15,
    phone = $16,
    is_pass_code = $17,
    is_active = $18
WHERE id = $1
RETURNING *;
