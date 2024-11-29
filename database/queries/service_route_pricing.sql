-- name: GetServiceRouteAvailabilityPricingById :one
SELECT * FROM service_route_pricing
WHERE id = $1
LIMIT 1;

-- name: GetServiceRouteAvailabilityPricingByServiceRouteAvailabilityIdAndFareTypeId :one
SELECT * FROM service_route_pricing
WHERE service_route_availability_id = $1 AND fare_type_id = $2
LIMIT 1;

-- name: CountServiceRouteAvailabilityPricesByServiceRouteAvailabilityId :one
SELECT
    COUNT(id)
FROM service_route_pricing
WHERE service_route_availability_id = $1;

-- name: ListServiceRouteAvailabilityPricesByServiceRouteAvailabilityId :many
SELECT * FROM service_route_pricing
WHERE service_route_availability_id = $1;

-- name: CountServiceRouteAvailabilityPricesByServiceRouteAvailabilityIdAndDateRange :one
SELECT
    COUNT(id)
FROM service_route_pricing
WHERE service_route_availability_id = $1 AND start_date <= $2 AND end_date >= $3;

-- name: ListServiceRouteAvailabilityPricesByServiceRouteAvailabilityIdAndDateRange :many
SELECT * FROM service_route_pricing
WHERE service_route_availability_id = $1 AND start_date <= $2 AND end_date >= $3;

-- name: CountServiceRouteAvailabilityPricesByServiceRouteAvailabilityIdAndDateRangeAndFareTypeId :one
SELECT
    COUNT(id)
FROM service_route_pricing
WHERE service_route_availability_id = $1 AND fare_type_id = $2 AND start_date <= $3 AND end_date >= $4;

-- name: ListServiceRouteAvailabilityPricesByServiceRouteAvailabilityIdAndDateRangeAndFareTypeId :many
SELECT * FROM service_route_pricing
WHERE service_route_availability_id = $1 AND fare_type_id = $2 AND start_date <= $3 AND end_date >= $4;

-- name: CreateServiceRouteAvailabilityPricing :one
INSERT INTO service_route_pricing (
    id,
    created,
    modified,
    service_route_availability_id,
    fare_type_id,
    price_type,
    adult_price,
    child_price,
    infant_price,
    start_date,
    end_date,
    excluded_dates,
    sold_out_dates,
    specific_terms,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
         )
RETURNING *;

-- name: UpsertServiceRouteAvailabilityPricing :one
INSERT INTO service_route_pricing (
    id,
    created,
    modified,
    service_route_availability_id,
    fare_type_id,
    price_type,
    adult_price,
    child_price,
    infant_price,
    start_date,
    end_date,
    excluded_dates,
    sold_out_dates,
    specific_terms,
    is_active
) VALUES (
             coalesce(nullif(@id, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
         )
ON CONFLICT (id) DO UPDATE
    SET
        modified = now() AT TIME ZONE 'UTC',
        service_route_availability_id = $1,
        fare_type_id = $2,
        price_type = $3,
        adult_price = $4,
        child_price = $5,
        infant_price = $6,
        start_date = $7,
        end_date = $8,
        excluded_dates = $9,
        sold_out_dates = $10,
        specific_terms = $11,
        is_active = $12
RETURNING *;

-- name: DeleteServiceRouteAvailabilityPricing :exec
DELETE FROM service_route_pricing
WHERE id = $1;

-- name: DeleteServiceRouteAvailabilityPricingByServiceRouteAvailabilityId :exec
DELETE FROM service_route_pricing
WHERE service_route_availability_id = $1;

-- name: UpdateServiceRouteAvailabilityPricing :one
UPDATE service_route_pricing
SET
    modified = now() AT TIME ZONE 'UTC',
    service_route_availability_id = $2,
    fare_type_id = $3,
    price_type = $4,
    adult_price = $5,
    child_price = $6,
    infant_price = $7,
    start_date = $8,
    end_date = $9,
    excluded_dates = $10,
    sold_out_dates = $11,
    specific_terms = $12,
    is_active = $13
WHERE id = $1
RETURNING *;
