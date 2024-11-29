-- name: GetCommissionById :one
SELECT * FROM commission
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetCommissionByCode :one
SELECT * FROM commission
WHERE code = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: GetCommissionByOperatorId :one
SELECT * FROM commission
WHERE operator_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: CountCommissions :one
SELECT
    COUNT(id)
FROM commission
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListCommissionsAsc :many
SELECT * FROM commission
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListCommissionsDesc :many
SELECT * FROM commission
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateCommission :one
INSERT INTO commission (
    id,
    created,
    modified,
    deleted,
    code,
    purpose,
    operator_id,
    type,
    value,
    description,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7
         )
RETURNING *;

-- name: DeleteCommission :exec
UPDATE commission
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateCommission :one
UPDATE commission
SET
    modified = now() AT TIME ZONE 'UTC',
    code = $2,
    purpose = $3,
    operator_id = $4,
    type = $5,
    value = $6,
    description = $7,
    is_active = $8
WHERE id = $1
RETURNING *;
