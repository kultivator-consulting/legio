-- name: GetServiceTermsConditionsById :one
SELECT * FROM service_terms_conditions
WHERE id = $1
LIMIT 1;

-- name: GetServiceTermsConditionsByServiceId :one
SELECT * FROM service_terms_conditions
WHERE service_id = $1
LIMIT 1;

-- name: CountServiceTermsConditionsByServiceId :one
SELECT
    COUNT(id)
FROM service_terms_conditions
WHERE service_id = $1;

-- name: ListServiceTermsConditionsByServiceId :many
SELECT * FROM service_terms_conditions
WHERE service_id = $1;

-- name: CreateServiceTermsConditions :one
INSERT INTO service_terms_conditions (
    id,
    created,
    modified,
    service_id,
    ordering,
    general_terms,
    luggage_notes,
    additional_notes
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5
         )
RETURNING *;

-- name: DeleteServiceTermsConditions :exec
DELETE FROM service_terms_conditions
WHERE id = $1;

-- name: DeleteServiceTermsConditionsByServiceId :exec
DELETE FROM service_terms_conditions
WHERE service_id = $1;

-- name: UpdateServiceTermsConditions :one
UPDATE service_terms_conditions
SET
    modified = now() AT TIME ZONE 'UTC',
    service_id = $2,
    ordering = $3,
    general_terms = $4,
    luggage_notes = $5,
    additional_notes = $6
WHERE id = $1
RETURNING *;
