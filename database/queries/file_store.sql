-- name: GetFileById :one
SELECT * FROM file_store
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPublicFileByFilename :one
SELECT * FROM file_store
WHERE filename = $1
  AND secure = false
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetFileByFilename :one
SELECT * FROM file_store
WHERE filename = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountFiles :one
SELECT
    COUNT(id)
FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListStoredFiles :many
SELECT stored_filename FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListBrowsableFilesAsc :many
SELECT * FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND browsable = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListBrowsableFilesDesc :many
SELECT * FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND browsable = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListFilesAsc :many
SELECT * FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListFilesDesc :many
SELECT * FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: UploadFile :one
INSERT INTO file_store (
    id,
    created,
    modified,
    deleted,
    is_system,
    filename,
    stored_filename,
    content_type,
    file_size,
    browsable,
    secure
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7
         )
RETURNING id;

-- name: ProcessingComplete :exec
UPDATE file_store
SET
    completed_processing = true
WHERE id = $1;

-- name: DeleteFile :exec
UPDATE file_store
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: DeleteFileByFilename :exec
UPDATE file_store
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE filename = $1;

-- name: UpdateFile :one
UPDATE file_store
SET
    modified = now() AT TIME ZONE 'UTC',
    is_system = $2,
    filename = $3,
    content_type = $4,
    file_size = $5,
    browsable = $6,
    secure = $7,
    completed_processing = false
WHERE id = $1
RETURNING *;

-- name: ProcessingFileComplete :one
UPDATE file_store
SET
    modified = now() AT TIME ZONE 'UTC',
    completed_processing = true
WHERE id = $1
RETURNING *;
