-- name: ListFileTagByFileStoreId :many
SELECT * FROM file_tag
WHERE file_store_id = $1;

-- name: DeleteFileTagsByFileStoreId :exec
DELETE FROM file_tag
WHERE file_store_id = $1;

-- name: AddFileTag :one
INSERT INTO file_tag (
    id,
    created,
    file_store_id,
    tag
) VALUES (
     uuid_generate_v4(),
     now() AT TIME ZONE 'UTC',
     $1, $2
)
RETURNING *;

-- name: DeleteFileTag :exec
DELETE FROM file_tag
WHERE id = $1;
