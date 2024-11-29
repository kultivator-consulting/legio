// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: pass.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countPasses = `-- name: CountPasses :one
SELECT
    COUNT(id)
FROM pass
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountPasses(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countPasses)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countPassesAtRoot = `-- name: CountPassesAtRoot :one
SELECT
    COUNT(id)
FROM pass
WHERE parent_pass_id IS NULL
  AND deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountPassesAtRoot(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countPassesAtRoot)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countPassesByParentPassId = `-- name: CountPassesByParentPassId :one
SELECT
    COUNT(id)
FROM pass
WHERE parent_pass_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountPassesByParentPassId(ctx context.Context, parentPassID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countPassesByParentPassId, parentPassID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPass = `-- name: CreatePass :one
INSERT INTO pass (
    id,
    created,
    modified,
    deleted,
    parent_pass_id,
    name,
    code,
    operator_codes,
    duration,
    duration_type,
    image,
    image_info,
    description,
    is_popular,
    is_top_up,
    ordering,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
         )
RETURNING id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active
`

type CreatePassParams struct {
	ParentPassID  pgtype.UUID `db:"parent_pass_id" json:"parentPassId"`
	Name          string      `db:"name" json:"name"`
	Code          string      `db:"code" json:"code"`
	OperatorCodes []string    `db:"operator_codes" json:"operatorCodes"`
	Duration      int32       `db:"duration" json:"duration"`
	DurationType  string      `db:"duration_type" json:"durationType"`
	Image         string      `db:"image" json:"image"`
	ImageInfo     string      `db:"image_info" json:"imageInfo"`
	Description   string      `db:"description" json:"description"`
	IsPopular     bool        `db:"is_popular" json:"isPopular"`
	IsTopUp       bool        `db:"is_top_up" json:"isTopUp"`
	Ordering      int32       `db:"ordering" json:"ordering"`
	IsActive      bool        `db:"is_active" json:"isActive"`
}

func (q *Queries) CreatePass(ctx context.Context, arg CreatePassParams) (Pass, error) {
	row := q.db.QueryRow(ctx, createPass,
		arg.ParentPassID,
		arg.Name,
		arg.Code,
		arg.OperatorCodes,
		arg.Duration,
		arg.DurationType,
		arg.Image,
		arg.ImageInfo,
		arg.Description,
		arg.IsPopular,
		arg.IsTopUp,
		arg.Ordering,
		arg.IsActive,
	)
	var i Pass
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.ParentPassID,
		&i.Name,
		&i.Code,
		&i.OperatorCodes,
		&i.Duration,
		&i.DurationType,
		&i.Image,
		&i.ImageInfo,
		&i.Description,
		&i.IsPopular,
		&i.IsTopUp,
		&i.Ordering,
		&i.IsActive,
	)
	return i, err
}

const deletePass = `-- name: DeletePass :exec
UPDATE pass
SET deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
`

func (q *Queries) DeletePass(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePass, id)
	return err
}

const getChildPassesByParentPassId = `-- name: GetChildPassesByParentPassId :many
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE parent_pass_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) GetChildPassesByParentPassId(ctx context.Context, parentPassID pgtype.UUID) ([]Pass, error) {
	rows, err := q.db.Query(ctx, getChildPassesByParentPassId, parentPassID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pass
	for rows.Next() {
		var i Pass
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.ParentPassID,
			&i.Name,
			&i.Code,
			&i.OperatorCodes,
			&i.Duration,
			&i.DurationType,
			&i.Image,
			&i.ImageInfo,
			&i.Description,
			&i.IsPopular,
			&i.IsTopUp,
			&i.Ordering,
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

const getPassByCode = `-- name: GetPassByCode :one
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE code = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetPassByCode(ctx context.Context, code string) (Pass, error) {
	row := q.db.QueryRow(ctx, getPassByCode, code)
	var i Pass
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.ParentPassID,
		&i.Name,
		&i.Code,
		&i.OperatorCodes,
		&i.Duration,
		&i.DurationType,
		&i.Image,
		&i.ImageInfo,
		&i.Description,
		&i.IsPopular,
		&i.IsTopUp,
		&i.Ordering,
		&i.IsActive,
	)
	return i, err
}

const getPassById = `-- name: GetPassById :one
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetPassById(ctx context.Context, id pgtype.UUID) (Pass, error) {
	row := q.db.QueryRow(ctx, getPassById, id)
	var i Pass
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.ParentPassID,
		&i.Name,
		&i.Code,
		&i.OperatorCodes,
		&i.Duration,
		&i.DurationType,
		&i.Image,
		&i.ImageInfo,
		&i.Description,
		&i.IsPopular,
		&i.IsTopUp,
		&i.Ordering,
		&i.IsActive,
	)
	return i, err
}

const getPassByName = `-- name: GetPassByName :one
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetPassByName(ctx context.Context, name string) (Pass, error) {
	row := q.db.QueryRow(ctx, getPassByName, name)
	var i Pass
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.ParentPassID,
		&i.Name,
		&i.Code,
		&i.OperatorCodes,
		&i.Duration,
		&i.DurationType,
		&i.Image,
		&i.ImageInfo,
		&i.Description,
		&i.IsPopular,
		&i.IsTopUp,
		&i.Ordering,
		&i.IsActive,
	)
	return i, err
}

const listPassesAsc = `-- name: ListPassesAsc :many
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListPassesAscParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPassesAsc(ctx context.Context, arg ListPassesAscParams) ([]Pass, error) {
	rows, err := q.db.Query(ctx, listPassesAsc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pass
	for rows.Next() {
		var i Pass
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.ParentPassID,
			&i.Name,
			&i.Code,
			&i.OperatorCodes,
			&i.Duration,
			&i.DurationType,
			&i.Image,
			&i.ImageInfo,
			&i.Description,
			&i.IsPopular,
			&i.IsTopUp,
			&i.Ordering,
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

const listPassesAtRootAsc = `-- name: ListPassesAtRootAsc :many
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE parent_pass_id IS NULL
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListPassesAtRootAscParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPassesAtRootAsc(ctx context.Context, arg ListPassesAtRootAscParams) ([]Pass, error) {
	rows, err := q.db.Query(ctx, listPassesAtRootAsc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pass
	for rows.Next() {
		var i Pass
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.ParentPassID,
			&i.Name,
			&i.Code,
			&i.OperatorCodes,
			&i.Duration,
			&i.DurationType,
			&i.Image,
			&i.ImageInfo,
			&i.Description,
			&i.IsPopular,
			&i.IsTopUp,
			&i.Ordering,
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

const listPassesAtRootDesc = `-- name: ListPassesAtRootDesc :many
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE parent_pass_id IS NULL
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text DESC
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListPassesAtRootDescParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPassesAtRootDesc(ctx context.Context, arg ListPassesAtRootDescParams) ([]Pass, error) {
	rows, err := q.db.Query(ctx, listPassesAtRootDesc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pass
	for rows.Next() {
		var i Pass
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.ParentPassID,
			&i.Name,
			&i.Code,
			&i.OperatorCodes,
			&i.Duration,
			&i.DurationType,
			&i.Image,
			&i.ImageInfo,
			&i.Description,
			&i.IsPopular,
			&i.IsTopUp,
			&i.Ordering,
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

const listPassesByParentPassIdAsc = `-- name: ListPassesByParentPassIdAsc :many
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE parent_pass_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY $2::text
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPassesByParentPassIdAscParams struct {
	ParentPassID      pgtype.UUID `db:"parent_pass_id" json:"parentPassId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPassesByParentPassIdAsc(ctx context.Context, arg ListPassesByParentPassIdAscParams) ([]Pass, error) {
	rows, err := q.db.Query(ctx, listPassesByParentPassIdAsc,
		arg.ParentPassID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pass
	for rows.Next() {
		var i Pass
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.ParentPassID,
			&i.Name,
			&i.Code,
			&i.OperatorCodes,
			&i.Duration,
			&i.DurationType,
			&i.Image,
			&i.ImageInfo,
			&i.Description,
			&i.IsPopular,
			&i.IsTopUp,
			&i.Ordering,
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

const listPassesByParentPassIdDesc = `-- name: ListPassesByParentPassIdDesc :many
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE parent_pass_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY $2::text DESC
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPassesByParentPassIdDescParams struct {
	ParentPassID      pgtype.UUID `db:"parent_pass_id" json:"parentPassId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPassesByParentPassIdDesc(ctx context.Context, arg ListPassesByParentPassIdDescParams) ([]Pass, error) {
	rows, err := q.db.Query(ctx, listPassesByParentPassIdDesc,
		arg.ParentPassID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pass
	for rows.Next() {
		var i Pass
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.ParentPassID,
			&i.Name,
			&i.Code,
			&i.OperatorCodes,
			&i.Duration,
			&i.DurationType,
			&i.Image,
			&i.ImageInfo,
			&i.Description,
			&i.IsPopular,
			&i.IsTopUp,
			&i.Ordering,
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

const listPassesDesc = `-- name: ListPassesDesc :many
SELECT id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active FROM pass
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text DESC
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListPassesDescParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPassesDesc(ctx context.Context, arg ListPassesDescParams) ([]Pass, error) {
	rows, err := q.db.Query(ctx, listPassesDesc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Pass
	for rows.Next() {
		var i Pass
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.ParentPassID,
			&i.Name,
			&i.Code,
			&i.OperatorCodes,
			&i.Duration,
			&i.DurationType,
			&i.Image,
			&i.ImageInfo,
			&i.Description,
			&i.IsPopular,
			&i.IsTopUp,
			&i.Ordering,
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

const updatePass = `-- name: UpdatePass :one
UPDATE pass
SET
    modified = now() AT TIME ZONE 'UTC',
    parent_pass_id = $2,
    name = $3,
    code = $4,
    operator_codes = $5,
    duration = $6,
    duration_type = $7,
    image = $8,
    image_info = $9,
    description = $10,
    is_popular = $11,
    is_top_up = $12,
    ordering = $13,
    is_active = $14
WHERE id = $1 AND deleted > now() AT TIME ZONE 'UTC'
RETURNING id, created, modified, deleted, parent_pass_id, name, code, operator_codes, duration, duration_type, image, image_info, description, is_popular, is_top_up, ordering, is_active
`

type UpdatePassParams struct {
	ID            pgtype.UUID `db:"id" json:"id"`
	ParentPassID  pgtype.UUID `db:"parent_pass_id" json:"parentPassId"`
	Name          string      `db:"name" json:"name"`
	Code          string      `db:"code" json:"code"`
	OperatorCodes []string    `db:"operator_codes" json:"operatorCodes"`
	Duration      int32       `db:"duration" json:"duration"`
	DurationType  string      `db:"duration_type" json:"durationType"`
	Image         string      `db:"image" json:"image"`
	ImageInfo     string      `db:"image_info" json:"imageInfo"`
	Description   string      `db:"description" json:"description"`
	IsPopular     bool        `db:"is_popular" json:"isPopular"`
	IsTopUp       bool        `db:"is_top_up" json:"isTopUp"`
	Ordering      int32       `db:"ordering" json:"ordering"`
	IsActive      bool        `db:"is_active" json:"isActive"`
}

func (q *Queries) UpdatePass(ctx context.Context, arg UpdatePassParams) (Pass, error) {
	row := q.db.QueryRow(ctx, updatePass,
		arg.ID,
		arg.ParentPassID,
		arg.Name,
		arg.Code,
		arg.OperatorCodes,
		arg.Duration,
		arg.DurationType,
		arg.Image,
		arg.ImageInfo,
		arg.Description,
		arg.IsPopular,
		arg.IsTopUp,
		arg.Ordering,
		arg.IsActive,
	)
	var i Pass
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.ParentPassID,
		&i.Name,
		&i.Code,
		&i.OperatorCodes,
		&i.Duration,
		&i.DurationType,
		&i.Image,
		&i.ImageInfo,
		&i.Description,
		&i.IsPopular,
		&i.IsTopUp,
		&i.Ordering,
		&i.IsActive,
	)
	return i, err
}