-- name: GetPagePathExtensionById :one
SELECT * FROM page_path_extension
WHERE id = $1
LIMIT 1;

-- name: CountPagePathExtensionsByPagePathId :one
SELECT
    COUNT(id)
FROM page_path_extension
WHERE page_path_id = $1;

-- name: ListPagePathExtensionsByPagePathId :many
SELECT * FROM page_path_extension
WHERE page_path_id = $1;

-- name: CreatePagePathExtension :one
INSERT INTO page_path_extension (
    id,
    created,
    modified,
    page_path_id,
    extension_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: DeletePagePathExtension :exec
DELETE FROM page_path_extension
WHERE id = $1;

-- name: DeletePagePathExtensionByPagePathId :exec
DELETE FROM page_path_extension
WHERE page_path_id = $1;

-- name: UpdatePagePathExtension :one
UPDATE page_path_extension
SET
    modified = now() AT TIME ZONE 'UTC',
    page_path_id = $2,
    extension_id = $3
WHERE id = $1
RETURNING *;
