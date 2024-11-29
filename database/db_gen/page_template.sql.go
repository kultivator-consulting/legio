// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: page_template.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countPageTemplates = `-- name: CountPageTemplates :one
SELECT
    COUNT(id)
FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountPageTemplates(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countPageTemplates)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPageTemplateAndReturnId = `-- name: CreatePageTemplateAndReturnId :one
INSERT INTO page_template (
    id,
    created,
    modified,
    deleted,
    domain_id,
    account_id,
    content_id,
    parent_page_path_id,
    title,
    slug,
    description,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8
         )
RETURNING id
`

type CreatePageTemplateAndReturnIdParams struct {
	DomainID         pgtype.UUID `db:"domain_id" json:"domainId"`
	AccountID        pgtype.UUID `db:"account_id" json:"accountId"`
	ContentID        pgtype.UUID `db:"content_id" json:"contentId"`
	ParentPagePathID pgtype.UUID `db:"parent_page_path_id" json:"parentPagePathId"`
	Title            string      `db:"title" json:"title"`
	Slug             string      `db:"slug" json:"slug"`
	Description      string      `db:"description" json:"description"`
	IsActive         bool        `db:"is_active" json:"isActive"`
}

func (q *Queries) CreatePageTemplateAndReturnId(ctx context.Context, arg CreatePageTemplateAndReturnIdParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, createPageTemplateAndReturnId,
		arg.DomainID,
		arg.AccountID,
		arg.ContentID,
		arg.ParentPagePathID,
		arg.Title,
		arg.Slug,
		arg.Description,
		arg.IsActive,
	)
	var id pgtype.UUID
	err := row.Scan(&id)
	return id, err
}

const deletePageTemplateById = `-- name: DeletePageTemplateById :exec
UPDATE page_template
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
`

func (q *Queries) DeletePageTemplateById(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePageTemplateById, id)
	return err
}

const getPageTemplateAtRootAndBySlug = `-- name: GetPageTemplateAtRootAndBySlug :one
SELECT id, created, modified, deleted, domain_id, account_id, content_id, parent_page_path_id, title, slug, description, is_active FROM page_template
WHERE parent_page_path_id IS NULL
  AND slug = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetPageTemplateAtRootAndBySlug(ctx context.Context, slug string) (PageTemplate, error) {
	row := q.db.QueryRow(ctx, getPageTemplateAtRootAndBySlug, slug)
	var i PageTemplate
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.ParentPagePathID,
		&i.Title,
		&i.Slug,
		&i.Description,
		&i.IsActive,
	)
	return i, err
}

const getPageTemplateById = `-- name: GetPageTemplateById :one
SELECT id, created, modified, deleted, domain_id, account_id, content_id, parent_page_path_id, title, slug, description, is_active FROM page_template
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetPageTemplateById(ctx context.Context, id pgtype.UUID) (PageTemplate, error) {
	row := q.db.QueryRow(ctx, getPageTemplateById, id)
	var i PageTemplate
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.ParentPagePathID,
		&i.Title,
		&i.Slug,
		&i.Description,
		&i.IsActive,
	)
	return i, err
}

const getPageTemplateByPagePathIdAndSlug = `-- name: GetPageTemplateByPagePathIdAndSlug :one
SELECT id, created, modified, deleted, domain_id, account_id, content_id, parent_page_path_id, title, slug, description, is_active FROM page_template
WHERE parent_page_path_id = $1
  AND slug = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

type GetPageTemplateByPagePathIdAndSlugParams struct {
	ParentPagePathID pgtype.UUID `db:"parent_page_path_id" json:"parentPagePathId"`
	Slug             string      `db:"slug" json:"slug"`
}

func (q *Queries) GetPageTemplateByPagePathIdAndSlug(ctx context.Context, arg GetPageTemplateByPagePathIdAndSlugParams) (PageTemplate, error) {
	row := q.db.QueryRow(ctx, getPageTemplateByPagePathIdAndSlug, arg.ParentPagePathID, arg.Slug)
	var i PageTemplate
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.ParentPagePathID,
		&i.Title,
		&i.Slug,
		&i.Description,
		&i.IsActive,
	)
	return i, err
}

const listAllPageTemplatesByPagePathIdAsc = `-- name: ListAllPageTemplatesByPagePathIdAsc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, parent_page_path_id, title, slug, description, is_active FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND parent_page_path_id = $1
ORDER BY title
`

func (q *Queries) ListAllPageTemplatesByPagePathIdAsc(ctx context.Context, parentPagePathID pgtype.UUID) ([]PageTemplate, error) {
	rows, err := q.db.Query(ctx, listAllPageTemplatesByPagePathIdAsc, parentPagePathID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PageTemplate
	for rows.Next() {
		var i PageTemplate
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.ParentPagePathID,
			&i.Title,
			&i.Slug,
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

const listPageTemplatesByDomainIdAndAtRootAsc = `-- name: ListPageTemplatesByDomainIdAndAtRootAsc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, parent_page_path_id, title, slug, description, is_active FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id IS NULL
ORDER BY $2::text
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPageTemplatesByDomainIdAndAtRootAscParams struct {
	DomainID          pgtype.UUID `db:"domain_id" json:"domainId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPageTemplatesByDomainIdAndAtRootAsc(ctx context.Context, arg ListPageTemplatesByDomainIdAndAtRootAscParams) ([]PageTemplate, error) {
	rows, err := q.db.Query(ctx, listPageTemplatesByDomainIdAndAtRootAsc,
		arg.DomainID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PageTemplate
	for rows.Next() {
		var i PageTemplate
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.ParentPagePathID,
			&i.Title,
			&i.Slug,
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

const listPageTemplatesByDomainIdAndAtRootDesc = `-- name: ListPageTemplatesByDomainIdAndAtRootDesc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, parent_page_path_id, title, slug, description, is_active FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id IS NULL
ORDER BY $2::text DESC
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPageTemplatesByDomainIdAndAtRootDescParams struct {
	DomainID          pgtype.UUID `db:"domain_id" json:"domainId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPageTemplatesByDomainIdAndAtRootDesc(ctx context.Context, arg ListPageTemplatesByDomainIdAndAtRootDescParams) ([]PageTemplate, error) {
	rows, err := q.db.Query(ctx, listPageTemplatesByDomainIdAndAtRootDesc,
		arg.DomainID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PageTemplate
	for rows.Next() {
		var i PageTemplate
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.ParentPagePathID,
			&i.Title,
			&i.Slug,
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

const listPageTemplatesByDomainIdAndPagePathIdAsc = `-- name: ListPageTemplatesByDomainIdAndPagePathIdAsc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, parent_page_path_id, title, slug, description, is_active FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id = $2
ORDER BY $3::text
OFFSET ($4::int - 1) * $5::int
    FETCH NEXT $5 ROWS ONLY
`

type ListPageTemplatesByDomainIdAndPagePathIdAscParams struct {
	DomainID          pgtype.UUID `db:"domain_id" json:"domainId"`
	ParentPagePathID  pgtype.UUID `db:"parent_page_path_id" json:"parentPagePathId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPageTemplatesByDomainIdAndPagePathIdAsc(ctx context.Context, arg ListPageTemplatesByDomainIdAndPagePathIdAscParams) ([]PageTemplate, error) {
	rows, err := q.db.Query(ctx, listPageTemplatesByDomainIdAndPagePathIdAsc,
		arg.DomainID,
		arg.ParentPagePathID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PageTemplate
	for rows.Next() {
		var i PageTemplate
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.ParentPagePathID,
			&i.Title,
			&i.Slug,
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

const listPageTemplatesByDomainIdAndPagePathIdDesc = `-- name: ListPageTemplatesByDomainIdAndPagePathIdDesc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, parent_page_path_id, title, slug, description, is_active FROM page_template
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id = $2
ORDER BY $3::text DESC
OFFSET ($4::int - 1) * $5::int
    FETCH NEXT $5 ROWS ONLY
`

type ListPageTemplatesByDomainIdAndPagePathIdDescParams struct {
	DomainID          pgtype.UUID `db:"domain_id" json:"domainId"`
	ParentPagePathID  pgtype.UUID `db:"parent_page_path_id" json:"parentPagePathId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPageTemplatesByDomainIdAndPagePathIdDesc(ctx context.Context, arg ListPageTemplatesByDomainIdAndPagePathIdDescParams) ([]PageTemplate, error) {
	rows, err := q.db.Query(ctx, listPageTemplatesByDomainIdAndPagePathIdDesc,
		arg.DomainID,
		arg.ParentPagePathID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PageTemplate
	for rows.Next() {
		var i PageTemplate
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.ParentPagePathID,
			&i.Title,
			&i.Slug,
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

const updatePageTemplateById = `-- name: UpdatePageTemplateById :one
UPDATE page_template
SET
    modified = now() AT TIME ZONE 'UTC',
    domain_id = $2,
    account_id = $3,
    content_id = $4,
    parent_page_path_id = $5,
    title = $6,
    slug = $7,
    description = $8,
    is_active = $9
WHERE id = $1
RETURNING id, created, modified, deleted, domain_id, account_id, content_id, parent_page_path_id, title, slug, description, is_active
`

type UpdatePageTemplateByIdParams struct {
	ID               pgtype.UUID `db:"id" json:"id"`
	DomainID         pgtype.UUID `db:"domain_id" json:"domainId"`
	AccountID        pgtype.UUID `db:"account_id" json:"accountId"`
	ContentID        pgtype.UUID `db:"content_id" json:"contentId"`
	ParentPagePathID pgtype.UUID `db:"parent_page_path_id" json:"parentPagePathId"`
	Title            string      `db:"title" json:"title"`
	Slug             string      `db:"slug" json:"slug"`
	Description      string      `db:"description" json:"description"`
	IsActive         bool        `db:"is_active" json:"isActive"`
}

func (q *Queries) UpdatePageTemplateById(ctx context.Context, arg UpdatePageTemplateByIdParams) (PageTemplate, error) {
	row := q.db.QueryRow(ctx, updatePageTemplateById,
		arg.ID,
		arg.DomainID,
		arg.AccountID,
		arg.ContentID,
		arg.ParentPagePathID,
		arg.Title,
		arg.Slug,
		arg.Description,
		arg.IsActive,
	)
	var i PageTemplate
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.ParentPagePathID,
		&i.Title,
		&i.Slug,
		&i.Description,
		&i.IsActive,
	)
	return i, err
}