-- name: GetPagePathById :one
SELECT * FROM page_path
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPagePathByDomainIdRootPathAndSlug :one
SELECT * FROM page_path
WHERE domain_id = $1
  AND slug = $2
  AND parent_page_path_id IS NULL
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPagePathByDomainIdAndSlug :one
SELECT * FROM page_path
WHERE domain_id = $1
  AND slug = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPagePathByDomainIdParentPagePathIdAndSlug :one
SELECT * FROM page_path
WHERE domain_id = $1
  AND slug = $2
  AND parent_page_path_id = $3
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountPagePaths :one
SELECT
    COUNT(id)
FROM page_path
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListAllPagePathsByParentPagePathIdAsc :many
SELECT * FROM page_path
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND parent_page_path_id = $1
ORDER BY title;

-- name: ListPagePathsByDomainIdAndParentPagePathIdAsc :many
SELECT * FROM page_path
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND domain_id = $1
    AND parent_page_path_id = $2
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPagePathsByDomainIdAndParentPagePathIdDesc :many
SELECT * FROM page_path
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id = $2
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPagePathsByDomainIdAtRootAsc :many
SELECT * FROM page_path
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND domain_id = $1
    AND parent_page_path_id IS NULL
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPagePathsByDomainIdAtRootDesc :many
SELECT * FROM page_path
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND domain_id = $1
    AND parent_page_path_id IS NULL
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListLinksAtRoot :many
SELECT * FROM page_path
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id IS NULL
ORDER BY 'title'::text;

-- name: ListLinksByParentId :many
SELECT * FROM page_path
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND domain_id = $1
  AND parent_page_path_id = $2
ORDER BY 'title'::text;

-- name: CreatePagePathAndReturnId :one
INSERT INTO page_path (
    id,
    created,
    modified,
    deleted,
    domain_id,
    account_id,
    parent_page_path_id,
    title,
    slug,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6
         )
RETURNING id;

-- name: DeletePagePath :exec
UPDATE page_path
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdatePagePath :one
UPDATE page_path
SET
    modified = now() AT TIME ZONE 'UTC',
    domain_id = $2,
    account_id = $3,
    parent_page_path_id = $4,
    title = $5,
    slug = $6,
    is_active = $7
WHERE id = $1
RETURNING *;
