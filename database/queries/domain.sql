-- name: GetDomain :one
SELECT * FROM domain
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetDomainById :one
SELECT * FROM domain
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetUnlockedDomainById :one
SELECT * FROM domain
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = FALSE
LIMIT 1;

-- name: GetDomainByDomainName :one
SELECT * FROM domain
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: CountDomains :one
SELECT
    COUNT(id)
FROM domain
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListDomainsAsc :many
SELECT * FROM domain
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListDomainsDesc :many
SELECT * FROM domain
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateDomain :one
INSERT INTO domain (
    id,
    created,
    modified,
    deleted,
    name,
    description,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3
         )
RETURNING *;

-- name: CreateDomainAndReturnId :one
INSERT INTO domain (
    id,
    created,
    modified,
    deleted,
    name,
    description,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3
         )
RETURNING id;

-- name: DeleteDomain :exec
UPDATE domain
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateDomain :one
UPDATE domain
SET
    modified = now() AT TIME ZONE 'UTC',
    name = $2,
    description = $3,
    is_active = $4
WHERE id = $1
RETURNING *;
