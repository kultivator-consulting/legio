-- name: GetPassAvailabilityPricingById :one
SELECT * FROM pass_availability_pricing
WHERE id = $1
LIMIT 1;

-- name: CountPassAvailabilityPricing :one
SELECT
    COUNT(id)
FROM pass_availability_pricing
WHERE pass_availability_id = $1;

-- name: ListPassAvailabilityPricingByPassAvailabilityId :many
SELECT * FROM pass_availability_pricing
WHERE pass_availability_id = $1;

-- name: CreatePassAvailabilityPricing :one
INSERT INTO pass_availability_pricing (
    id,
    created,
    modified,
    pass_availability_id,
    fare_type_id,
    adult_price,
    child_price,
    infant_price,
    start_date,
    end_date,
    excluded_dates
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8
         )
RETURNING *;

-- name: UpsertPassAvailabilityPricing :one
INSERT INTO pass_availability_pricing (
    id,
    created,
    modified,
    pass_availability_id,
    fare_type_id,
    adult_price,
    child_price,
    infant_price,
    start_date,
    end_date,
    excluded_dates
) VALUES (
             coalesce(nullif(@id, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
                $1, $2, $3, $4, $5, $6, $7, $8
         )
ON CONFLICT (id) DO UPDATE
    SET
        modified = now() AT TIME ZONE 'UTC',
        pass_availability_id = $1,
        fare_type_id = $2,
        adult_price = $3,
        child_price = $4,
        infant_price = $5,
        start_date = $6,
        end_date = $7,
        excluded_dates = $8
RETURNING *;

-- name: DeletePassAvailabilityPricing :exec
DELETE FROM pass_availability_pricing
WHERE id = $1;

-- name: UpdatePassAvailabilityPricing :one
UPDATE pass_availability_pricing
SET
    modified = now() AT TIME ZONE 'UTC',
    pass_availability_id = $2,
    fare_type_id = $3,
    adult_price = $4,
    child_price = $5,
    infant_price = $6,
    start_date = $7,
    end_date = $8,
    excluded_dates = $9
WHERE id = $1
RETURNING *;
