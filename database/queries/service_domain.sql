-- name: GetServiceDomainById :one
SELECT * FROM service_domain
WHERE id = $1
LIMIT 1;

-- name: GetServiceDomainByDomainId :one
SELECT * FROM service_domain
WHERE domain_id = $1
LIMIT 1;

-- name: CountServiceDomainsByServiceId :one
SELECT
    COUNT(id)
FROM service_domain
WHERE service_id = $1;

-- name: ListServiceDomainsByServiceId :many
SELECT * FROM service_domain
WHERE service_id = $1;

-- name: CountServiceDomainsByDomainId :one
SELECT
    COUNT(id)
FROM service_domain
WHERE domain_id = $1;

-- name: ListServiceDomainsByDomainId :many
SELECT * FROM service_domain
WHERE domain_id = $1;

-- name: CreateServiceDomain :one
INSERT INTO service_domain (
    id,
    created,
    modified,
    service_id,
    domain_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: DeleteServiceDomain :exec
DELETE FROM service_domain
WHERE id = $1;

-- name: DeleteServiceDomainByServiceId :exec
DELETE FROM service_domain
WHERE service_id = $1;

-- name: UpdateServiceDomain :one
UPDATE service_domain
SET
    modified = now() AT TIME ZONE 'UTC',
    service_id = $2,
    domain_id = $3
WHERE id = $1
RETURNING *;
