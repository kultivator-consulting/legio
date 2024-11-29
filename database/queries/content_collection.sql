-- name: GetContentCollectionById :one
SELECT * FROM content_collection
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetContentCollectionByParentIdAndContentId :one
SELECT * FROM content_collection
WHERE parent_id = $1
  AND content_id = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: ListContentCollectionByParentId :many
SELECT * FROM content_collection
WHERE parent_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY ordering;

-- name: CreateContentCollectionAndReturnId :one
INSERT INTO content_collection (
    id,
    created,
    modified,
    deleted,
    parent_id,
    content_id,
    ordering,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING id;

-- name: DeleteContentCollectionByIdAndParentId :exec
UPDATE content_collection
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
  AND parent_id = $2;

-- name: DeleteContentCollection :exec
UPDATE content_collection
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateContentCollectionById :one
UPDATE content_collection
SET
    modified = now() AT TIME ZONE 'UTC',
    parent_id = $2,
    content_id = $3,
    ordering = $4,
    is_active = $5
WHERE id = $1
RETURNING *;
