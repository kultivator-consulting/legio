// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: extension.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countExtensions = `-- name: CountExtensions :one
SELECT
    COUNT(id)
FROM extension
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountExtensions(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countExtensions)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createExtension = `-- name: CreateExtension :one
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
RETURNING id, created, modified, deleted, name, slug, icon, data, is_active
`

type CreateExtensionParams struct {
	Name     string `db:"name" json:"name"`
	Slug     string `db:"slug" json:"slug"`
	Icon     string `db:"icon" json:"icon"`
	Data     string `db:"data" json:"data"`
	IsActive bool   `db:"is_active" json:"isActive"`
}

func (q *Queries) CreateExtension(ctx context.Context, arg CreateExtensionParams) (Extension, error) {
	row := q.db.QueryRow(ctx, createExtension,
		arg.Name,
		arg.Slug,
		arg.Icon,
		arg.Data,
		arg.IsActive,
	)
	var i Extension
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Name,
		&i.Slug,
		&i.Icon,
		&i.Data,
		&i.IsActive,
	)
	return i, err
}

const deleteExtension = `-- name: DeleteExtension :exec
UPDATE extension
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
`

func (q *Queries) DeleteExtension(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteExtension, id)
	return err
}

const getExtensionById = `-- name: GetExtensionById :one
SELECT id, created, modified, deleted, name, slug, icon, data, is_active FROM extension
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetExtensionById(ctx context.Context, id pgtype.UUID) (Extension, error) {
	row := q.db.QueryRow(ctx, getExtensionById, id)
	var i Extension
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Name,
		&i.Slug,
		&i.Icon,
		&i.Data,
		&i.IsActive,
	)
	return i, err
}

const listExtensions = `-- name: ListExtensions :many
SELECT id, created, modified, deleted, name, slug, icon, data, is_active FROM extension
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) ListExtensions(ctx context.Context) ([]Extension, error) {
	rows, err := q.db.Query(ctx, listExtensions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Extension
	for rows.Next() {
		var i Extension
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.Name,
			&i.Slug,
			&i.Icon,
			&i.Data,
			&i.IsActive,
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

const updateExtension = `-- name: UpdateExtension :one
UPDATE extension
SET
    modified = now() AT TIME ZONE 'UTC',
    name = $2,
    slug = $3,
    icon = $4,
    data = $5,
    is_active = $6
WHERE id = $1
RETURNING id, created, modified, deleted, name, slug, icon, data, is_active
`

type UpdateExtensionParams struct {
	ID       pgtype.UUID `db:"id" json:"id"`
	Name     string      `db:"name" json:"name"`
	Slug     string      `db:"slug" json:"slug"`
	Icon     string      `db:"icon" json:"icon"`
	Data     string      `db:"data" json:"data"`
	IsActive bool        `db:"is_active" json:"isActive"`
}

func (q *Queries) UpdateExtension(ctx context.Context, arg UpdateExtensionParams) (Extension, error) {
	row := q.db.QueryRow(ctx, updateExtension,
		arg.ID,
		arg.Name,
		arg.Slug,
		arg.Icon,
		arg.Data,
		arg.IsActive,
	)
	var i Extension
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Name,
		&i.Slug,
		&i.Icon,
		&i.Data,
		&i.IsActive,
	)
	return i, err
}
