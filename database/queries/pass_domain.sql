-- name: GetPassDomainById :one
SELECT * FROM pass_domain
WHERE id = $1
LIMIT 1;

-- name: CountPassDomains :one
SELECT
    COUNT(id)
FROM pass_domain
WHERE pass_id = $1;

-- name: ListPassDomainsByPassId :many
SELECT * FROM pass_domain
WHERE pass_id = $1;

-- name: CreatePassDomain :one
INSERT INTO pass_domain (
    id,
    created,
    modified,
    pass_id,
    domain_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: DeletePassDomain :exec
DELETE FROM pass_domain
WHERE id = $1;

-- name: DeletePassDomainByPassId :exec
DELETE FROM pass_domain
WHERE pass_id = $1;

-- name: UpdatePassDomain :one
UPDATE pass_domain
SET
    modified = now() AT TIME ZONE 'UTC',
    pass_id = $2,
    domain_id = $3
WHERE id = $1
RETURNING *;
