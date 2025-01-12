// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: commission.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countCommissions = `-- name: CountCommissions :one
SELECT
    COUNT(id)
FROM commission
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountCommissions(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countCommissions)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createCommission = `-- name: CreateCommission :one
INSERT INTO commission (
    id,
    created,
    modified,
    deleted,
    code,
    purpose,
    operator_id,
    type,
    value,
    description,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7
         )
RETURNING id, created, modified, deleted, code, purpose, operator_id, type, value, description, is_active
`

type CreateCommissionParams struct {
	Code        string         `db:"code" json:"code"`
	Purpose     string         `db:"purpose" json:"purpose"`
	OperatorID  pgtype.UUID    `db:"operator_id" json:"operatorId"`
	Type        string         `db:"type" json:"type"`
	Value       pgtype.Numeric `db:"value" json:"value"`
	Description pgtype.Text    `db:"description" json:"description"`
	IsActive    bool           `db:"is_active" json:"isActive"`
}

func (q *Queries) CreateCommission(ctx context.Context, arg CreateCommissionParams) (Commission, error) {
	row := q.db.QueryRow(ctx, createCommission,
		arg.Code,
		arg.Purpose,
		arg.OperatorID,
		arg.Type,
		arg.Value,
		arg.Description,
		arg.IsActive,
	)
	var i Commission
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Code,
		&i.Purpose,
		&i.OperatorID,
		&i.Type,
		&i.Value,
		&i.Description,
		&i.IsActive,
	)
	return i, err
}

const deleteCommission = `-- name: DeleteCommission :exec
UPDATE commission
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
`

func (q *Queries) DeleteCommission(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteCommission, id)
	return err
}

const getCommissionByCode = `-- name: GetCommissionByCode :one
SELECT id, created, modified, deleted, code, purpose, operator_id, type, value, description, is_active FROM commission
WHERE code = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1
`

func (q *Queries) GetCommissionByCode(ctx context.Context, code string) (Commission, error) {
	row := q.db.QueryRow(ctx, getCommissionByCode, code)
	var i Commission
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Code,
		&i.Purpose,
		&i.OperatorID,
		&i.Type,
		&i.Value,
		&i.Description,
		&i.IsActive,
	)
	return i, err
}

const getCommissionById = `-- name: GetCommissionById :one
SELECT id, created, modified, deleted, code, purpose, operator_id, type, value, description, is_active FROM commission
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetCommissionById(ctx context.Context, id pgtype.UUID) (Commission, error) {
	row := q.db.QueryRow(ctx, getCommissionById, id)
	var i Commission
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Code,
		&i.Purpose,
		&i.OperatorID,
		&i.Type,
		&i.Value,
		&i.Description,
		&i.IsActive,
	)
	return i, err
}

const getCommissionByOperatorId = `-- name: GetCommissionByOperatorId :one
SELECT id, created, modified, deleted, code, purpose, operator_id, type, value, description, is_active FROM commission
WHERE operator_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1
`

func (q *Queries) GetCommissionByOperatorId(ctx context.Context, operatorID pgtype.UUID) (Commission, error) {
	row := q.db.QueryRow(ctx, getCommissionByOperatorId, operatorID)
	var i Commission
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Code,
		&i.Purpose,
		&i.OperatorID,
		&i.Type,
		&i.Value,
		&i.Description,
		&i.IsActive,
	)
	return i, err
}

const listCommissionsAsc = `-- name: ListCommissionsAsc :many
SELECT id, created, modified, deleted, code, purpose, operator_id, type, value, description, is_active FROM commission
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListCommissionsAscParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListCommissionsAsc(ctx context.Context, arg ListCommissionsAscParams) ([]Commission, error) {
	rows, err := q.db.Query(ctx, listCommissionsAsc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Commission
	for rows.Next() {
		var i Commission
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.Code,
			&i.Purpose,
			&i.OperatorID,
			&i.Type,
			&i.Value,
			&i.Description,
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

const listCommissionsDesc = `-- name: ListCommissionsDesc :many
SELECT id, created, modified, deleted, code, purpose, operator_id, type, value, description, is_active FROM commission
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text DESC
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListCommissionsDescParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListCommissionsDesc(ctx context.Context, arg ListCommissionsDescParams) ([]Commission, error) {
	rows, err := q.db.Query(ctx, listCommissionsDesc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Commission
	for rows.Next() {
		var i Commission
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.Code,
			&i.Purpose,
			&i.OperatorID,
			&i.Type,
			&i.Value,
			&i.Description,
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

const updateCommission = `-- name: UpdateCommission :one
UPDATE commission
SET
    modified = now() AT TIME ZONE 'UTC',
    code = $2,
    purpose = $3,
    operator_id = $4,
    type = $5,
    value = $6,
    description = $7,
    is_active = $8
WHERE id = $1
RETURNING id, created, modified, deleted, code, purpose, operator_id, type, value, description, is_active
`

type UpdateCommissionParams struct {
	ID          pgtype.UUID    `db:"id" json:"id"`
	Code        string         `db:"code" json:"code"`
	Purpose     string         `db:"purpose" json:"purpose"`
	OperatorID  pgtype.UUID    `db:"operator_id" json:"operatorId"`
	Type        string         `db:"type" json:"type"`
	Value       pgtype.Numeric `db:"value" json:"value"`
	Description pgtype.Text    `db:"description" json:"description"`
	IsActive    bool           `db:"is_active" json:"isActive"`
}

func (q *Queries) UpdateCommission(ctx context.Context, arg UpdateCommissionParams) (Commission, error) {
	row := q.db.QueryRow(ctx, updateCommission,
		arg.ID,
		arg.Code,
		arg.Purpose,
		arg.OperatorID,
		arg.Type,
		arg.Value,
		arg.Description,
		arg.IsActive,
	)
	var i Commission
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Code,
		&i.Purpose,
		&i.OperatorID,
		&i.Type,
		&i.Value,
		&i.Description,
		&i.IsActive,
	)
	return i, err
}
