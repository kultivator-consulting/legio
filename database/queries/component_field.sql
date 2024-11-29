-- name: GetComponentFieldById :one
SELECT * FROM component_field
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: ListComponentFieldByComponentId :many
SELECT * FROM component_field
WHERE component_id = $1
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: ListComponentFieldByComponentIdAndName :many
SELECT * FROM component_field
WHERE component_id = $1
  AND name = $2
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: CreateComponentFieldAndReturnId :one
INSERT INTO component_field (
    id,
    created,
    modified,
    deleted,
    component_id,
    name,
    description,
    data_type,
    editor_type,
    validation,
    default_value,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8
         )
RETURNING id;

-- name: DeleteComponentFieldByIdAndComponentId :exec
UPDATE component_field
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
    AND component_id = $2;

-- name: DeleteComponentFieldByComponentId :exec
UPDATE component_field
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE component_id = $1;

-- name: UpdateComponentField :one
UPDATE component_field
SET
    modified = now() AT TIME ZONE 'UTC',
    component_id = $2,
    name = $3,
    description = $4,
    data_type = $5,
    editor_type = $6,
    validation = $7,
    default_value = $8,
    is_active = $9
WHERE id = $1
RETURNING *;
