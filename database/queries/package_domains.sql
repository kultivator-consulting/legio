-- name: GetPackageDomainById :one
SELECT * FROM package_domains
WHERE id = $1
LIMIT 1;

-- name: GetPackageDomainsByPackageId :many
SELECT * FROM package_domains
WHERE package_id = $1;

-- name: CountPackageDomainsByPackageId :one
SELECT
    COUNT(id)
FROM package_domains
WHERE package_id = $1;

-- name: ListPackageDomainsByDomainIdAsc :many
SELECT * FROM package_domains
WHERE domain_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageDomainsByDomainIdDesc :many
SELECT * FROM package_domains
WHERE domain_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageDomainsByPackageIdAsc :many
SELECT * FROM package_domains
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageDomainsByPackageIdDesc :many
SELECT * FROM package_domains
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackageDomain :one
INSERT INTO package_domains (
    id,
    created,
    modified,
    package_id,
    domain_id
) VALUES (
     uuid_generate_v4(),
     now() AT TIME ZONE 'UTC',
     now() AT TIME ZONE 'UTC',
     $1, $2
 )
RETURNING *;

-- name: UpdatePackageDomain :one
UPDATE package_domains
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    domain_id = $3
WHERE id = $1
RETURNING *;

-- name: DeletePackageDomain :exec
DELETE FROM package_domains
WHERE id = $1;

-- name: DeletePackageDomainsByPackageId :exec
DELETE FROM package_domains
WHERE package_id = $1;