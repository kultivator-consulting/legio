-- name: GetFileMetaByKey :one
SELECT * FROM file_meta
WHERE key = $1 LIMIT 1;

-- name: ListFileMetaByFileStoreId :many
SELECT * FROM file_meta
WHERE file_store_id = $1;

-- name: ListAttachmentsByFileStoreId :many
SELECT attached_file_store_id FROM file_meta
WHERE file_store_id = $1;

-- name: DeleteFileMetaByFileStoreId :exec
DELETE FROM file_meta
WHERE file_store_id = $1;

-- name: AddFileMeta :one
INSERT INTO file_meta (
    id,
    created,
    file_store_id,
    key,
    value,
    attached_file_store_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING *;
