// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: file_store.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countFiles = `-- name: CountFiles :one
SELECT
    COUNT(id)
FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountFiles(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countFiles)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const deleteFile = `-- name: DeleteFile :exec
UPDATE file_store
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
`

func (q *Queries) DeleteFile(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteFile, id)
	return err
}

const deleteFileByFilename = `-- name: DeleteFileByFilename :exec
UPDATE file_store
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE filename = $1
`

func (q *Queries) DeleteFileByFilename(ctx context.Context, filename string) error {
	_, err := q.db.Exec(ctx, deleteFileByFilename, filename)
	return err
}

const getFileByFilename = `-- name: GetFileByFilename :one
SELECT id, created, modified, deleted, is_system, filename, stored_filename, content_type, file_size, browsable, secure, completed_processing FROM file_store
WHERE filename = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetFileByFilename(ctx context.Context, filename string) (FileStore, error) {
	row := q.db.QueryRow(ctx, getFileByFilename, filename)
	var i FileStore
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.IsSystem,
		&i.Filename,
		&i.StoredFilename,
		&i.ContentType,
		&i.FileSize,
		&i.Browsable,
		&i.Secure,
		&i.CompletedProcessing,
	)
	return i, err
}

const getFileById = `-- name: GetFileById :one
SELECT id, created, modified, deleted, is_system, filename, stored_filename, content_type, file_size, browsable, secure, completed_processing FROM file_store
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetFileById(ctx context.Context, id pgtype.UUID) (FileStore, error) {
	row := q.db.QueryRow(ctx, getFileById, id)
	var i FileStore
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.IsSystem,
		&i.Filename,
		&i.StoredFilename,
		&i.ContentType,
		&i.FileSize,
		&i.Browsable,
		&i.Secure,
		&i.CompletedProcessing,
	)
	return i, err
}

const getPublicFileByFilename = `-- name: GetPublicFileByFilename :one
SELECT id, created, modified, deleted, is_system, filename, stored_filename, content_type, file_size, browsable, secure, completed_processing FROM file_store
WHERE filename = $1
  AND secure = false
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetPublicFileByFilename(ctx context.Context, filename string) (FileStore, error) {
	row := q.db.QueryRow(ctx, getPublicFileByFilename, filename)
	var i FileStore
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.IsSystem,
		&i.Filename,
		&i.StoredFilename,
		&i.ContentType,
		&i.FileSize,
		&i.Browsable,
		&i.Secure,
		&i.CompletedProcessing,
	)
	return i, err
}

const listBrowsableFilesAsc = `-- name: ListBrowsableFilesAsc :many
SELECT id, created, modified, deleted, is_system, filename, stored_filename, content_type, file_size, browsable, secure, completed_processing FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND browsable = $1
ORDER BY $2::text
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListBrowsableFilesAscParams struct {
	Browsable         bool   `db:"browsable" json:"browsable"`
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListBrowsableFilesAsc(ctx context.Context, arg ListBrowsableFilesAscParams) ([]FileStore, error) {
	rows, err := q.db.Query(ctx, listBrowsableFilesAsc,
		arg.Browsable,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FileStore
	for rows.Next() {
		var i FileStore
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.IsSystem,
			&i.Filename,
			&i.StoredFilename,
			&i.ContentType,
			&i.FileSize,
			&i.Browsable,
			&i.Secure,
			&i.CompletedProcessing,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBrowsableFilesDesc = `-- name: ListBrowsableFilesDesc :many
SELECT id, created, modified, deleted, is_system, filename, stored_filename, content_type, file_size, browsable, secure, completed_processing FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND browsable = $1
ORDER BY $2::text DESC
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListBrowsableFilesDescParams struct {
	Browsable         bool   `db:"browsable" json:"browsable"`
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListBrowsableFilesDesc(ctx context.Context, arg ListBrowsableFilesDescParams) ([]FileStore, error) {
	rows, err := q.db.Query(ctx, listBrowsableFilesDesc,
		arg.Browsable,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FileStore
	for rows.Next() {
		var i FileStore
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.IsSystem,
			&i.Filename,
			&i.StoredFilename,
			&i.ContentType,
			&i.FileSize,
			&i.Browsable,
			&i.Secure,
			&i.CompletedProcessing,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listFilesAsc = `-- name: ListFilesAsc :many
SELECT id, created, modified, deleted, is_system, filename, stored_filename, content_type, file_size, browsable, secure, completed_processing FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListFilesAscParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListFilesAsc(ctx context.Context, arg ListFilesAscParams) ([]FileStore, error) {
	rows, err := q.db.Query(ctx, listFilesAsc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FileStore
	for rows.Next() {
		var i FileStore
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.IsSystem,
			&i.Filename,
			&i.StoredFilename,
			&i.ContentType,
			&i.FileSize,
			&i.Browsable,
			&i.Secure,
			&i.CompletedProcessing,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listFilesDesc = `-- name: ListFilesDesc :many
SELECT id, created, modified, deleted, is_system, filename, stored_filename, content_type, file_size, browsable, secure, completed_processing FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text DESC
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListFilesDescParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListFilesDesc(ctx context.Context, arg ListFilesDescParams) ([]FileStore, error) {
	rows, err := q.db.Query(ctx, listFilesDesc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FileStore
	for rows.Next() {
		var i FileStore
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.IsSystem,
			&i.Filename,
			&i.StoredFilename,
			&i.ContentType,
			&i.FileSize,
			&i.Browsable,
			&i.Secure,
			&i.CompletedProcessing,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listStoredFiles = `-- name: ListStoredFiles :many
SELECT stored_filename FROM file_store
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) ListStoredFiles(ctx context.Context) ([]string, error) {
	rows, err := q.db.Query(ctx, listStoredFiles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var stored_filename string
		if err := rows.Scan(&stored_filename); err != nil {
			return nil, err
		}
		items = append(items, stored_filename)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const processingComplete = `-- name: ProcessingComplete :exec
UPDATE file_store
SET
    completed_processing = true
WHERE id = $1
`

func (q *Queries) ProcessingComplete(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, processingComplete, id)
	return err
}

const processingFileComplete = `-- name: ProcessingFileComplete :one
UPDATE file_store
SET
    modified = now() AT TIME ZONE 'UTC',
    completed_processing = true
WHERE id = $1
RETURNING id, created, modified, deleted, is_system, filename, stored_filename, content_type, file_size, browsable, secure, completed_processing
`

func (q *Queries) ProcessingFileComplete(ctx context.Context, id pgtype.UUID) (FileStore, error) {
	row := q.db.QueryRow(ctx, processingFileComplete, id)
	var i FileStore
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.IsSystem,
		&i.Filename,
		&i.StoredFilename,
		&i.ContentType,
		&i.FileSize,
		&i.Browsable,
		&i.Secure,
		&i.CompletedProcessing,
	)
	return i, err
}

const updateFile = `-- name: UpdateFile :one
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
RETURNING id, created, modified, deleted, is_system, filename, stored_filename, content_type, file_size, browsable, secure, completed_processing
`

type UpdateFileParams struct {
	ID          pgtype.UUID `db:"id" json:"id"`
	IsSystem    bool        `db:"is_system" json:"isSystem"`
	Filename    string      `db:"filename" json:"filename"`
	ContentType pgtype.Text `db:"content_type" json:"contentType"`
	FileSize    int64       `db:"file_size" json:"fileSize"`
	Browsable   bool        `db:"browsable" json:"browsable"`
	Secure      bool        `db:"secure" json:"secure"`
}

func (q *Queries) UpdateFile(ctx context.Context, arg UpdateFileParams) (FileStore, error) {
	row := q.db.QueryRow(ctx, updateFile,
		arg.ID,
		arg.IsSystem,
		arg.Filename,
		arg.ContentType,
		arg.FileSize,
		arg.Browsable,
		arg.Secure,
	)
	var i FileStore
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.IsSystem,
		&i.Filename,
		&i.StoredFilename,
		&i.ContentType,
		&i.FileSize,
		&i.Browsable,
		&i.Secure,
		&i.CompletedProcessing,
	)
	return i, err
}

const uploadFile = `-- name: UploadFile :one
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
RETURNING id
`

type UploadFileParams struct {
	IsSystem       bool        `db:"is_system" json:"isSystem"`
	Filename       string      `db:"filename" json:"filename"`
	StoredFilename string      `db:"stored_filename" json:"storedFilename"`
	ContentType    pgtype.Text `db:"content_type" json:"contentType"`
	FileSize       int64       `db:"file_size" json:"fileSize"`
	Browsable      bool        `db:"browsable" json:"browsable"`
	Secure         bool        `db:"secure" json:"secure"`
}

func (q *Queries) UploadFile(ctx context.Context, arg UploadFileParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, uploadFile,
		arg.IsSystem,
		arg.Filename,
		arg.StoredFilename,
		arg.ContentType,
		arg.FileSize,
		arg.Browsable,
		arg.Secure,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}
