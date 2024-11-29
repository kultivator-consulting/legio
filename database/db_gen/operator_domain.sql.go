// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: operator_domain.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countOperatorDomainsByOperatorId = `-- name: CountOperatorDomainsByOperatorId :one
SELECT
    COUNT(id)
FROM operator_domain
WHERE operator_id = $1
`

func (q *Queries) CountOperatorDomainsByOperatorId(ctx context.Context, operatorID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countOperatorDomainsByOperatorId, operatorID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createOperatorDomain = `-- name: CreateOperatorDomain :one
INSERT INTO operator_domain (
    id,
    created,
    modified,
    operator_id,
    domain_id
) VALUES (
    uuid_generate_v4(),
    now() AT TIME ZONE 'UTC',
    now() AT TIME ZONE 'UTC',
    $1, $2
)
RETURNING id, created, modified, operator_id, domain_id
`

type CreateOperatorDomainParams struct {
	OperatorID pgtype.UUID `db:"operator_id" json:"operatorId"`
	DomainID   pgtype.UUID `db:"domain_id" json:"domainId"`
}

func (q *Queries) CreateOperatorDomain(ctx context.Context, arg CreateOperatorDomainParams) (OperatorDomain, error) {
	row := q.db.QueryRow(ctx, createOperatorDomain, arg.OperatorID, arg.DomainID)
	var i OperatorDomain
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.OperatorID,
		&i.DomainID,
	)
	return i, err
}

const deleteOperatorDomain = `-- name: DeleteOperatorDomain :exec
DELETE FROM operator_domain
WHERE id = $1
`

func (q *Queries) DeleteOperatorDomain(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteOperatorDomain, id)
	return err
}

const deleteOperatorDomainsByOperatorId = `-- name: DeleteOperatorDomainsByOperatorId :exec
DELETE FROM operator_domain
WHERE operator_id = $1
`

func (q *Queries) DeleteOperatorDomainsByOperatorId(ctx context.Context, operatorID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteOperatorDomainsByOperatorId, operatorID)
	return err
}

const getOperatorDomainById = `-- name: GetOperatorDomainById :one
SELECT id, created, modified, operator_id, domain_id FROM operator_domain
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetOperatorDomainById(ctx context.Context, id pgtype.UUID) (OperatorDomain, error) {
	row := q.db.QueryRow(ctx, getOperatorDomainById, id)
	var i OperatorDomain
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.OperatorID,
		&i.DomainID,
	)
	return i, err
}

const getProductByName = `-- name: GetProductByName :one
SELECT id, created, modified, deleted, product_type_id, operator_id, start_location_id, end_location_id, name, start_place, end_place, operator_code, instructions, notes, is_active FROM product
WHERE name = $1
LIMIT 1
`

func (q *Queries) GetProductByName(ctx context.Context, name string) (Product, error) {
	row := q.db.QueryRow(ctx, getProductByName, name)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.ProductTypeID,
		&i.OperatorID,
		&i.StartLocationID,
		&i.EndLocationID,
		&i.Name,
		&i.StartPlace,
		&i.EndPlace,
		&i.OperatorCode,
		&i.Instructions,
		&i.Notes,
		&i.IsActive,
	)
	return i, err
}

const listOperatorDomainsByOperatorId = `-- name: ListOperatorDomainsByOperatorId :many
SELECT id, created, modified, operator_id, domain_id FROM operator_domain
WHERE operator_id = $1
`

func (q *Queries) ListOperatorDomainsByOperatorId(ctx context.Context, operatorID pgtype.UUID) ([]OperatorDomain, error) {
	rows, err := q.db.Query(ctx, listOperatorDomainsByOperatorId, operatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []OperatorDomain
	for rows.Next() {
		var i OperatorDomain
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.OperatorID,
			&i.DomainID,
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

const updateOperatorDomain = `-- name: UpdateOperatorDomain :one
UPDATE operator_domain
SET
    modified = now() AT TIME ZONE 'UTC',
    operator_id = $2,
    domain_id = $3
WHERE id = $1
RETURNING id, created, modified, operator_id, domain_id
`

type UpdateOperatorDomainParams struct {
	ID         pgtype.UUID `db:"id" json:"id"`
	OperatorID pgtype.UUID `db:"operator_id" json:"operatorId"`
	DomainID   pgtype.UUID `db:"domain_id" json:"domainId"`
}

func (q *Queries) UpdateOperatorDomain(ctx context.Context, arg UpdateOperatorDomainParams) (OperatorDomain, error) {
	row := q.db.QueryRow(ctx, updateOperatorDomain, arg.ID, arg.OperatorID, arg.DomainID)
	var i OperatorDomain
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.OperatorID,
		&i.DomainID,
	)
	return i, err
}
