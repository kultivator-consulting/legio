// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: extra_product_types.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countExtraProductTypesByExtraId = `-- name: CountExtraProductTypesByExtraId :one
SELECT
    COUNT(id)
FROM extra_product_type
WHERE extra_id = $1
`

func (q *Queries) CountExtraProductTypesByExtraId(ctx context.Context, extraID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countExtraProductTypesByExtraId, extraID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countExtraProductTypesByProductTypeId = `-- name: CountExtraProductTypesByProductTypeId :one
SELECT
    COUNT(id)
FROM extra_product_type
WHERE product_type_id = $1
`

func (q *Queries) CountExtraProductTypesByProductTypeId(ctx context.Context, productTypeID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countExtraProductTypesByProductTypeId, productTypeID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createExtraProductType = `-- name: CreateExtraProductType :one
INSERT INTO extra_product_type (
    id,
    created,
    modified,
    extra_id,
    product_type_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING id, created, modified, extra_id, product_type_id
`

type CreateExtraProductTypeParams struct {
	ExtraID       pgtype.UUID `db:"extra_id" json:"extraId"`
	ProductTypeID pgtype.UUID `db:"product_type_id" json:"productTypeId"`
}

func (q *Queries) CreateExtraProductType(ctx context.Context, arg CreateExtraProductTypeParams) (ExtraProductType, error) {
	row := q.db.QueryRow(ctx, createExtraProductType, arg.ExtraID, arg.ProductTypeID)
	var i ExtraProductType
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ExtraID,
		&i.ProductTypeID,
	)
	return i, err
}

const deleteExtraProductType = `-- name: DeleteExtraProductType :exec
DELETE FROM extra_product_type
WHERE id = $1
`

func (q *Queries) DeleteExtraProductType(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteExtraProductType, id)
	return err
}

const deleteExtraProductTypeByExtraId = `-- name: DeleteExtraProductTypeByExtraId :exec
DELETE FROM extra_product_type
WHERE extra_id = $1
`

func (q *Queries) DeleteExtraProductTypeByExtraId(ctx context.Context, extraID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteExtraProductTypeByExtraId, extraID)
	return err
}

const getExtraProductTypeById = `-- name: GetExtraProductTypeById :one
SELECT id, created, modified, extra_id, product_type_id FROM extra_product_type
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetExtraProductTypeById(ctx context.Context, id pgtype.UUID) (ExtraProductType, error) {
	row := q.db.QueryRow(ctx, getExtraProductTypeById, id)
	var i ExtraProductType
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ExtraID,
		&i.ProductTypeID,
	)
	return i, err
}

const listExtraProductTypesByExtraId = `-- name: ListExtraProductTypesByExtraId :many
SELECT id, created, modified, extra_id, product_type_id FROM extra_product_type
WHERE extra_id = $1
`

func (q *Queries) ListExtraProductTypesByExtraId(ctx context.Context, extraID pgtype.UUID) ([]ExtraProductType, error) {
	rows, err := q.db.Query(ctx, listExtraProductTypesByExtraId, extraID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExtraProductType
	for rows.Next() {
		var i ExtraProductType
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.ExtraID,
			&i.ProductTypeID,
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

const listExtraProductTypesByProductTypeId = `-- name: ListExtraProductTypesByProductTypeId :many
SELECT id, created, modified, extra_id, product_type_id FROM extra_product_type
WHERE product_type_id = $1
`

func (q *Queries) ListExtraProductTypesByProductTypeId(ctx context.Context, productTypeID pgtype.UUID) ([]ExtraProductType, error) {
	rows, err := q.db.Query(ctx, listExtraProductTypesByProductTypeId, productTypeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ExtraProductType
	for rows.Next() {
		var i ExtraProductType
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.ExtraID,
			&i.ProductTypeID,
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

const updateExtraProductType = `-- name: UpdateExtraProductType :one
UPDATE extra_product_type
SET
    modified = now() AT TIME ZONE 'UTC',
    extra_id = $2,
    product_type_id = $3
WHERE id = $1
RETURNING id, created, modified, extra_id, product_type_id
`

type UpdateExtraProductTypeParams struct {
	ID            pgtype.UUID `db:"id" json:"id"`
	ExtraID       pgtype.UUID `db:"extra_id" json:"extraId"`
	ProductTypeID pgtype.UUID `db:"product_type_id" json:"productTypeId"`
}

func (q *Queries) UpdateExtraProductType(ctx context.Context, arg UpdateExtraProductTypeParams) (ExtraProductType, error) {
	row := q.db.QueryRow(ctx, updateExtraProductType, arg.ID, arg.ExtraID, arg.ProductTypeID)
	var i ExtraProductType
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ExtraID,
		&i.ProductTypeID,
	)
	return i, err
}
