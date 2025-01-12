// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: page.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countPages = `-- name: CountPages :one
SELECT
    COUNT(id)
FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountPages(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countPages)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPage = `-- name: CreatePage :one
INSERT INTO page (
    id,
    created,
    modified,
    deleted,
    domain_id,
    account_id,
    content_id,
    page_path_id,
    title,
    slug,
    seo_title,
    seo_description,
    seo_keywords,
    draft_page_id,
    page_template_id,
    publish_at,
    unpublish_at,
    version,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
         )
RETURNING id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active
`

type CreatePageParams struct {
	DomainID       pgtype.UUID        `db:"domain_id" json:"domainId"`
	AccountID      pgtype.UUID        `db:"account_id" json:"accountId"`
	ContentID      pgtype.UUID        `db:"content_id" json:"contentId"`
	PagePathID     pgtype.UUID        `db:"page_path_id" json:"pagePathId"`
	Title          string             `db:"title" json:"title"`
	Slug           string             `db:"slug" json:"slug"`
	SeoTitle       string             `db:"seo_title" json:"seoTitle"`
	SeoDescription string             `db:"seo_description" json:"seoDescription"`
	SeoKeywords    string             `db:"seo_keywords" json:"seoKeywords"`
	DraftPageID    pgtype.UUID        `db:"draft_page_id" json:"draftPageId"`
	PageTemplateID pgtype.UUID        `db:"page_template_id" json:"pageTemplateId"`
	PublishAt      pgtype.Timestamptz `db:"publish_at" json:"publishAt"`
	UnpublishAt    pgtype.Timestamptz `db:"unpublish_at" json:"unpublishAt"`
	Version        int32              `db:"version" json:"version"`
	IsActive       bool               `db:"is_active" json:"isActive"`
}

func (q *Queries) CreatePage(ctx context.Context, arg CreatePageParams) (Page, error) {
	row := q.db.QueryRow(ctx, createPage,
		arg.DomainID,
		arg.AccountID,
		arg.ContentID,
		arg.PagePathID,
		arg.Title,
		arg.Slug,
		arg.SeoTitle,
		arg.SeoDescription,
		arg.SeoKeywords,
		arg.DraftPageID,
		arg.PageTemplateID,
		arg.PublishAt,
		arg.UnpublishAt,
		arg.Version,
		arg.IsActive,
	)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.PagePathID,
		&i.Title,
		&i.Slug,
		&i.SeoTitle,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.DraftPageID,
		&i.PageTemplateID,
		&i.PublishAt,
		&i.UnpublishAt,
		&i.Version,
		&i.IsActive,
	)
	return i, err
}

const deletePageById = `-- name: DeletePageById :exec
UPDATE page
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
`

func (q *Queries) DeletePageById(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePageById, id)
	return err
}

const detachTemplateFromPageByTemplateId = `-- name: DetachTemplateFromPageByTemplateId :exec
UPDATE page
SET
    page_template_id = NULL
WHERE page_template_id = $1
`

func (q *Queries) DetachTemplateFromPageByTemplateId(ctx context.Context, pageTemplateID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, detachTemplateFromPageByTemplateId, pageTemplateID)
	return err
}

const getPageByDomainIdAndId = `-- name: GetPageByDomainIdAndId :one
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE domain_id = $1
  AND id = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

type GetPageByDomainIdAndIdParams struct {
	DomainID pgtype.UUID `db:"domain_id" json:"domainId"`
	ID       pgtype.UUID `db:"id" json:"id"`
}

func (q *Queries) GetPageByDomainIdAndId(ctx context.Context, arg GetPageByDomainIdAndIdParams) (Page, error) {
	row := q.db.QueryRow(ctx, getPageByDomainIdAndId, arg.DomainID, arg.ID)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.PagePathID,
		&i.Title,
		&i.Slug,
		&i.SeoTitle,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.DraftPageID,
		&i.PageTemplateID,
		&i.PublishAt,
		&i.UnpublishAt,
		&i.Version,
		&i.IsActive,
	)
	return i, err
}

const getPageByDomainIdPagePathIdAndSlug = `-- name: GetPageByDomainIdPagePathIdAndSlug :one
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE domain_id = $1
  AND slug = $2
  AND page_path_id = $3
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

type GetPageByDomainIdPagePathIdAndSlugParams struct {
	DomainID   pgtype.UUID `db:"domain_id" json:"domainId"`
	Slug       string      `db:"slug" json:"slug"`
	PagePathID pgtype.UUID `db:"page_path_id" json:"pagePathId"`
}

func (q *Queries) GetPageByDomainIdPagePathIdAndSlug(ctx context.Context, arg GetPageByDomainIdPagePathIdAndSlugParams) (Page, error) {
	row := q.db.QueryRow(ctx, getPageByDomainIdPagePathIdAndSlug, arg.DomainID, arg.Slug, arg.PagePathID)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.PagePathID,
		&i.Title,
		&i.Slug,
		&i.SeoTitle,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.DraftPageID,
		&i.PageTemplateID,
		&i.PublishAt,
		&i.UnpublishAt,
		&i.Version,
		&i.IsActive,
	)
	return i, err
}

const getPageByDomainIdRootPathAndSlug = `-- name: GetPageByDomainIdRootPathAndSlug :one
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE domain_id = $1
  AND slug = $2
  AND page_path_id IS NULL
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

type GetPageByDomainIdRootPathAndSlugParams struct {
	DomainID pgtype.UUID `db:"domain_id" json:"domainId"`
	Slug     string      `db:"slug" json:"slug"`
}

func (q *Queries) GetPageByDomainIdRootPathAndSlug(ctx context.Context, arg GetPageByDomainIdRootPathAndSlugParams) (Page, error) {
	row := q.db.QueryRow(ctx, getPageByDomainIdRootPathAndSlug, arg.DomainID, arg.Slug)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.PagePathID,
		&i.Title,
		&i.Slug,
		&i.SeoTitle,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.DraftPageID,
		&i.PageTemplateID,
		&i.PublishAt,
		&i.UnpublishAt,
		&i.Version,
		&i.IsActive,
	)
	return i, err
}

const getPageById = `-- name: GetPageById :one
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetPageById(ctx context.Context, id pgtype.UUID) (Page, error) {
	row := q.db.QueryRow(ctx, getPageById, id)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.PagePathID,
		&i.Title,
		&i.Slug,
		&i.SeoTitle,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.DraftPageID,
		&i.PageTemplateID,
		&i.PublishAt,
		&i.UnpublishAt,
		&i.Version,
		&i.IsActive,
	)
	return i, err
}

const getPageByPagePathIdAndSlug = `-- name: GetPageByPagePathIdAndSlug :one
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE page_path_id = $1
  AND slug = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

type GetPageByPagePathIdAndSlugParams struct {
	PagePathID pgtype.UUID `db:"page_path_id" json:"pagePathId"`
	Slug       string      `db:"slug" json:"slug"`
}

func (q *Queries) GetPageByPagePathIdAndSlug(ctx context.Context, arg GetPageByPagePathIdAndSlugParams) (Page, error) {
	row := q.db.QueryRow(ctx, getPageByPagePathIdAndSlug, arg.PagePathID, arg.Slug)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.PagePathID,
		&i.Title,
		&i.Slug,
		&i.SeoTitle,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.DraftPageID,
		&i.PageTemplateID,
		&i.PublishAt,
		&i.UnpublishAt,
		&i.Version,
		&i.IsActive,
	)
	return i, err
}

const listAllPagesByPagePathIdAsc = `-- name: ListAllPagesByPagePathIdAsc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND page_path_id = $1
ORDER BY title
`

func (q *Queries) ListAllPagesByPagePathIdAsc(ctx context.Context, pagePathID pgtype.UUID) ([]Page, error) {
	rows, err := q.db.Query(ctx, listAllPagesByPagePathIdAsc, pagePathID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.PagePathID,
			&i.Title,
			&i.Slug,
			&i.SeoTitle,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.DraftPageID,
			&i.PageTemplateID,
			&i.PublishAt,
			&i.UnpublishAt,
			&i.Version,
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

const listPagesAtRoot = `-- name: ListPagesAtRoot :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND page_path_id IS NULL
ORDER BY 'title'::text
`

func (q *Queries) ListPagesAtRoot(ctx context.Context) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesAtRoot)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.PagePathID,
			&i.Title,
			&i.Slug,
			&i.SeoTitle,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.DraftPageID,
			&i.PageTemplateID,
			&i.PublishAt,
			&i.UnpublishAt,
			&i.Version,
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

const listPagesByDomainIdAndPagePathIdAsc = `-- name: ListPagesByDomainIdAndPagePathIdAsc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND domain_id = $1
    AND page_path_id = $2
ORDER BY $3::text
OFFSET ($4::int - 1) * $5::int
    FETCH NEXT $5 ROWS ONLY
`

type ListPagesByDomainIdAndPagePathIdAscParams struct {
	DomainID          pgtype.UUID `db:"domain_id" json:"domainId"`
	PagePathID        pgtype.UUID `db:"page_path_id" json:"pagePathId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPagesByDomainIdAndPagePathIdAsc(ctx context.Context, arg ListPagesByDomainIdAndPagePathIdAscParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByDomainIdAndPagePathIdAsc,
		arg.DomainID,
		arg.PagePathID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.PagePathID,
			&i.Title,
			&i.Slug,
			&i.SeoTitle,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.DraftPageID,
			&i.PageTemplateID,
			&i.PublishAt,
			&i.UnpublishAt,
			&i.Version,
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

const listPagesByDomainIdAndPagePathIdDesc = `-- name: ListPagesByDomainIdAndPagePathIdDesc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND domain_id = $1
    AND page_path_id = $2
ORDER BY $3::text DESC
OFFSET ($4::int - 1) * $5::int
    FETCH NEXT $5 ROWS ONLY
`

type ListPagesByDomainIdAndPagePathIdDescParams struct {
	DomainID          pgtype.UUID `db:"domain_id" json:"domainId"`
	PagePathID        pgtype.UUID `db:"page_path_id" json:"pagePathId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPagesByDomainIdAndPagePathIdDesc(ctx context.Context, arg ListPagesByDomainIdAndPagePathIdDescParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByDomainIdAndPagePathIdDesc,
		arg.DomainID,
		arg.PagePathID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.PagePathID,
			&i.Title,
			&i.Slug,
			&i.SeoTitle,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.DraftPageID,
			&i.PageTemplateID,
			&i.PublishAt,
			&i.UnpublishAt,
			&i.Version,
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

const listPagesByDomainIdAsc = `-- name: ListPagesByDomainIdAsc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND domain_id = $1
ORDER BY $2::text
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPagesByDomainIdAscParams struct {
	DomainID          pgtype.UUID `db:"domain_id" json:"domainId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPagesByDomainIdAsc(ctx context.Context, arg ListPagesByDomainIdAscParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByDomainIdAsc,
		arg.DomainID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.PagePathID,
			&i.Title,
			&i.Slug,
			&i.SeoTitle,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.DraftPageID,
			&i.PageTemplateID,
			&i.PublishAt,
			&i.UnpublishAt,
			&i.Version,
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

const listPagesByDomainIdDesc = `-- name: ListPagesByDomainIdDesc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
ORDER BY $2::text DESC
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPagesByDomainIdDescParams struct {
	DomainID          pgtype.UUID `db:"domain_id" json:"domainId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPagesByDomainIdDesc(ctx context.Context, arg ListPagesByDomainIdDescParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByDomainIdDesc,
		arg.DomainID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.PagePathID,
			&i.Title,
			&i.Slug,
			&i.SeoTitle,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.DraftPageID,
			&i.PageTemplateID,
			&i.PublishAt,
			&i.UnpublishAt,
			&i.Version,
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

const listPagesByPagePathId = `-- name: ListPagesByPagePathId :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND page_path_id = $1
ORDER BY 'title'::text
`

func (q *Queries) ListPagesByPagePathId(ctx context.Context, pagePathID pgtype.UUID) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByPagePathId, pagePathID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.PagePathID,
			&i.Title,
			&i.Slug,
			&i.SeoTitle,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.DraftPageID,
			&i.PageTemplateID,
			&i.PublishAt,
			&i.UnpublishAt,
			&i.Version,
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

const listPagesByPagePathIdAsc = `-- name: ListPagesByPagePathIdAsc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND page_path_id = $1
ORDER BY $2::text
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPagesByPagePathIdAscParams struct {
	PagePathID        pgtype.UUID `db:"page_path_id" json:"pagePathId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPagesByPagePathIdAsc(ctx context.Context, arg ListPagesByPagePathIdAscParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByPagePathIdAsc,
		arg.PagePathID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.PagePathID,
			&i.Title,
			&i.Slug,
			&i.SeoTitle,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.DraftPageID,
			&i.PageTemplateID,
			&i.PublishAt,
			&i.UnpublishAt,
			&i.Version,
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

const listPagesByPagePathIdDesc = `-- name: ListPagesByPagePathIdDesc :many
SELECT id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND page_path_id = $1
ORDER BY $2::text DESC
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPagesByPagePathIdDescParams struct {
	PagePathID        pgtype.UUID `db:"page_path_id" json:"pagePathId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPagesByPagePathIdDesc(ctx context.Context, arg ListPagesByPagePathIdDescParams) ([]Page, error) {
	rows, err := q.db.Query(ctx, listPagesByPagePathIdDesc,
		arg.PagePathID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Page
	for rows.Next() {
		var i Page
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.DomainID,
			&i.AccountID,
			&i.ContentID,
			&i.PagePathID,
			&i.Title,
			&i.Slug,
			&i.SeoTitle,
			&i.SeoDescription,
			&i.SeoKeywords,
			&i.DraftPageID,
			&i.PageTemplateID,
			&i.PublishAt,
			&i.UnpublishAt,
			&i.Version,
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

const updatePageById = `-- name: UpdatePageById :one
UPDATE page
SET
    modified = now() AT TIME ZONE 'UTC',
    domain_id = $2,
    account_id = $3,
    content_id = $4,
    page_path_id = $5,
    title = $6,
    slug = $7,
    seo_title = $8,
    seo_description = $9,
    seo_keywords = $10,
    draft_page_id = $11,
    page_template_id = $12,
    publish_at = $13,
    unpublish_at = $14,
    version = $15,
    is_active = $16
WHERE id = $1
RETURNING id, created, modified, deleted, domain_id, account_id, content_id, page_path_id, title, slug, seo_title, seo_description, seo_keywords, draft_page_id, page_template_id, publish_at, unpublish_at, version, is_active
`

type UpdatePageByIdParams struct {
	ID             pgtype.UUID        `db:"id" json:"id"`
	DomainID       pgtype.UUID        `db:"domain_id" json:"domainId"`
	AccountID      pgtype.UUID        `db:"account_id" json:"accountId"`
	ContentID      pgtype.UUID        `db:"content_id" json:"contentId"`
	PagePathID     pgtype.UUID        `db:"page_path_id" json:"pagePathId"`
	Title          string             `db:"title" json:"title"`
	Slug           string             `db:"slug" json:"slug"`
	SeoTitle       string             `db:"seo_title" json:"seoTitle"`
	SeoDescription string             `db:"seo_description" json:"seoDescription"`
	SeoKeywords    string             `db:"seo_keywords" json:"seoKeywords"`
	DraftPageID    pgtype.UUID        `db:"draft_page_id" json:"draftPageId"`
	PageTemplateID pgtype.UUID        `db:"page_template_id" json:"pageTemplateId"`
	PublishAt      pgtype.Timestamptz `db:"publish_at" json:"publishAt"`
	UnpublishAt    pgtype.Timestamptz `db:"unpublish_at" json:"unpublishAt"`
	Version        int32              `db:"version" json:"version"`
	IsActive       bool               `db:"is_active" json:"isActive"`
}

func (q *Queries) UpdatePageById(ctx context.Context, arg UpdatePageByIdParams) (Page, error) {
	row := q.db.QueryRow(ctx, updatePageById,
		arg.ID,
		arg.DomainID,
		arg.AccountID,
		arg.ContentID,
		arg.PagePathID,
		arg.Title,
		arg.Slug,
		arg.SeoTitle,
		arg.SeoDescription,
		arg.SeoKeywords,
		arg.DraftPageID,
		arg.PageTemplateID,
		arg.PublishAt,
		arg.UnpublishAt,
		arg.Version,
		arg.IsActive,
	)
	var i Page
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.DomainID,
		&i.AccountID,
		&i.ContentID,
		&i.PagePathID,
		&i.Title,
		&i.Slug,
		&i.SeoTitle,
		&i.SeoDescription,
		&i.SeoKeywords,
		&i.DraftPageID,
		&i.PageTemplateID,
		&i.PublishAt,
		&i.UnpublishAt,
		&i.Version,
		&i.IsActive,
	)
	return i, err
}
