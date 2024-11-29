-- name: GetComponentById :one
SELECT * FROM component
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetComponentByName :one
SELECT * FROM component
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetComponentByClassName :one
SELECT * FROM component
WHERE class_name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetComponentByHtmlTag :one
SELECT * FROM component
WHERE html_tag = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountComponents :one
SELECT
    COUNT(id)
FROM component
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListComponentsAsc :many
SELECT * FROM component
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListComponentsDesc :many
SELECT * FROM component
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateComponentAndReturnId :one
INSERT INTO component (
    id,
    created,
    modified,
    deleted,
    name,
    description,
    icon,
    class_name,
    html_tag,
    child_tag_constraints,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7
         )
RETURNING id;

-- name: DeleteComponent :exec
UPDATE component
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateComponent :one
UPDATE component
SET
    modified = now() AT TIME ZONE 'UTC',
    name = $2,
    description = $3,
    icon = $4,
    class_name = $5,
    html_tag = $6,
    child_tag_constraints = $7,
    is_active = $8
WHERE id = $1
RETURNING *;
