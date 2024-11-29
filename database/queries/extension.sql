-- name: GetExtensionById :one
SELECT * FROM extension
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountExtensions :one
SELECT
    COUNT(id)
FROM extension
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListExtensions :many
SELECT * FROM extension
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: CreateExtension :one
INSERT INTO extension (
    id,
    created,
    modified,
    deleted,
    name,
    slug,
    icon,
    data,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5
         )
RETURNING *;

-- name: DeleteExtension :exec
UPDATE extension
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateExtension :one
UPDATE extension
SET
    modified = now() AT TIME ZONE 'UTC',
    name = $2,
    slug = $3,
    icon = $4,
    data = $5,
    is_active = $6
WHERE id = $1
RETURNING *;
