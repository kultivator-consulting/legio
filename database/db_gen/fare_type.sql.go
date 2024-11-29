// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: fare_type.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countFareTypes = `-- name: CountFareTypes :one
SELECT
    COUNT(id)
FROM fare_type
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1
`

func (q *Queries) CountFareTypes(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countFareTypes)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createFareType = `-- name: CreateFareType :one
INSERT INTO fare_type (
    id,
    created,
    modified,
    deleted,
    name,
    ordering,
    description,
    is_international,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5
         )
RETURNING id, created, modified, deleted, name, ordering, description, is_international, is_active
`

type CreateFareTypeParams struct {
	Name            string `db:"name" json:"name"`
	Ordering        int32  `db:"ordering" json:"ordering"`
	Description     string `db:"description" json:"description"`
	IsInternational bool   `db:"is_international" json:"isInternational"`
	IsActive        bool   `db:"is_active" json:"isActive"`
}

func (q *Queries) CreateFareType(ctx context.Context, arg CreateFareTypeParams) (FareType, error) {
	row := q.db.QueryRow(ctx, createFareType,
		arg.Name,
		arg.Ordering,
		arg.Description,
		arg.IsInternational,
		arg.IsActive,
	)
	var i FareType
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Name,
		&i.Ordering,
		&i.Description,
		&i.IsInternational,
		&i.IsActive,
	)
	return i, err
}

const deleteFareType = `-- name: DeleteFareType :exec
UPDATE fare_type
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
`

func (q *Queries) DeleteFareType(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteFareType, id)
	return err
}

const getActiveFareTypes = `-- name: GetActiveFareTypes :many
SELECT id, created, modified, deleted, name, ordering, description, is_international, is_active FROM fare_type
WHERE is_active = true
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY ordering
`

func (q *Queries) GetActiveFareTypes(ctx context.Context) ([]FareType, error) {
	rows, err := q.db.Query(ctx, getActiveFareTypes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FareType
	for rows.Next() {
		var i FareType
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.Name,
			&i.Ordering,
			&i.Description,
			&i.IsInternational,
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

const getFareTypeById = `-- name: GetFareTypeById :one
SELECT id, created, modified, deleted, name, ordering, description, is_international, is_active FROM fare_type
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetFareTypeById(ctx context.Context, id pgtype.UUID) (FareType, error) {
	row := q.db.QueryRow(ctx, getFareTypeById, id)
	var i FareType
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Name,
		&i.Ordering,
		&i.Description,
		&i.IsInternational,
		&i.IsActive,
	)
	return i, err
}

const getFareTypeByName = `-- name: GetFareTypeByName :one
SELECT id, created, modified, deleted, name, ordering, description, is_international, is_active FROM fare_type
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetFareTypeByName(ctx context.Context, name string) (FareType, error) {
	row := q.db.QueryRow(ctx, getFareTypeByName, name)
	var i FareType
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Name,
		&i.Ordering,
		&i.Description,
		&i.IsInternational,
		&i.IsActive,
	)
	return i, err
}

const listFareTypesAsc = `-- name: ListFareTypesAsc :many
SELECT id, created, modified, deleted, name, ordering, description, is_international, is_active FROM fare_type
    WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListFareTypesAscParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListFareTypesAsc(ctx context.Context, arg ListFareTypesAscParams) ([]FareType, error) {
	rows, err := q.db.Query(ctx, listFareTypesAsc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FareType
	for rows.Next() {
		var i FareType
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.Name,
			&i.Ordering,
			&i.Description,
			&i.IsInternational,
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

const listFareTypesDesc = `-- name: ListFareTypesDesc :many
SELECT id, created, modified, deleted, name, ordering, description, is_international, is_active FROM fare_type
    WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text DESC
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListFareTypesDescParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListFareTypesDesc(ctx context.Context, arg ListFareTypesDescParams) ([]FareType, error) {
	rows, err := q.db.Query(ctx, listFareTypesDesc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FareType
	for rows.Next() {
		var i FareType
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.Name,
			&i.Ordering,
			&i.Description,
			&i.IsInternational,
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

const updateFareType = `-- name: UpdateFareType :one
UPDATE fare_type
SET
    modified = now() AT TIME ZONE 'UTC',
    name = $2,
    ordering = $3,
    description = $4,
    is_international = $5,
    is_active = $6
WHERE id = $1
RETURNING id, created, modified, deleted, name, ordering, description, is_international, is_active
`

type UpdateFareTypeParams struct {
	ID              pgtype.UUID `db:"id" json:"id"`
	Name            string      `db:"name" json:"name"`
	Ordering        int32       `db:"ordering" json:"ordering"`
	Description     string      `db:"description" json:"description"`
	IsInternational bool        `db:"is_international" json:"isInternational"`
	IsActive        bool        `db:"is_active" json:"isActive"`
}

func (q *Queries) UpdateFareType(ctx context.Context, arg UpdateFareTypeParams) (FareType, error) {
	row := q.db.QueryRow(ctx, updateFareType,
		arg.ID,
		arg.Name,
		arg.Ordering,
		arg.Description,
		arg.IsInternational,
		arg.IsActive,
	)
	var i FareType
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Name,
		&i.Ordering,
		&i.Description,
		&i.IsInternational,
		&i.IsActive,
	)
	return i, err
}
