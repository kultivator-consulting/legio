-- name: GetPageById :one
SELECT * FROM page
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPageByDomainIdAndId :one
SELECT * FROM page
WHERE domain_id = $1
  AND id = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPageByPagePathIdAndSlug :one
SELECT * FROM page
WHERE page_path_id = $1
  AND slug = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPageByDomainIdRootPathAndSlug :one
SELECT * FROM page
WHERE domain_id = $1
  AND slug = $2
  AND page_path_id IS NULL
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPageByDomainIdPagePathIdAndSlug :one
SELECT * FROM page
WHERE domain_id = $1
  AND slug = $2
  AND page_path_id = $3
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountPages :one
SELECT
    COUNT(id)
FROM page
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListAllPagesByPagePathIdAsc :many
SELECT * FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND page_path_id = $1
ORDER BY title;

-- name: ListPagesByPagePathIdAsc :many
SELECT * FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND page_path_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPagesByPagePathIdDesc :many
SELECT * FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND page_path_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPagesByDomainIdAsc :many
SELECT * FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND domain_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPagesByDomainIdDesc :many
SELECT * FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPagesByDomainIdAndPagePathIdAsc :many
SELECT * FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND domain_id = $1
    AND page_path_id = $2
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPagesByDomainIdAndPagePathIdDesc :many
SELECT * FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND domain_id = $1
    AND page_path_id = $2
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPagesAtRoot :many
SELECT * FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND page_path_id IS NULL
ORDER BY 'title'::text;

-- name: ListPagesByPagePathId :many
SELECT * FROM page
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND page_path_id = $1
ORDER BY 'title'::text;

-- name: CreatePage :one
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
RETURNING *;

-- name: DeletePageById :exec
UPDATE page
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: DetachTemplateFromPageByTemplateId :exec
UPDATE page
SET
    page_template_id = NULL
WHERE page_template_id = $1;

-- name: UpdatePageById :one
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
RETURNING *;
