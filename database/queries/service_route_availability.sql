-- name: GetServiceRouteAvailabilityById :one
SELECT * FROM service_route_availability
WHERE id = $1
LIMIT 1;

-- name: GetServiceRouteAvailabilityByCiTimeId :one
SELECT * FROM service_route_availability
WHERE ci_time_id = $1
LIMIT 1;

-- name: CountServiceRouteAvailabilitiesByServiceRouteId :one
SELECT
    COUNT(id)
FROM service_route_availability
WHERE service_route_id = $1;

-- name: ListServiceRouteAvailabilitiesByServiceRouteId :many
SELECT * FROM service_route_availability
WHERE service_route_id = $1;

-- name: CountServiceRouteAvailabilitiesByServiceIdAndDateRange :one
SELECT
    COUNT(id)
FROM service_route_availability
WHERE service_route_id = $1 AND start_date <= $2 AND end_date >= $3;

-- name: ListServiceRouteAvailabilitiesByServiceIdAndDateRange :many
SELECT * FROM service_route_availability
WHERE service_route_id = $1 AND start_date <= $2 AND end_date >= $3;

-- name: ServiceRouteAvailabilitiesSearch :many
SELECT service.name,
       service.sub_service_name,
       service.description,
       service.excluded_dates AS service_excluded_dates,
       service.sold_out_dates AS service_sold_out_dates,
       service_route.start_station_id,
       service_route.end_station_id,
       service_route.is_active,
       service_route.is_main_route,
       service_route.ic_route_id,
       service_route_availability.start_date,
       service_route_availability.end_date,
       service_route_availability.departure_time,
       service_route_availability.arrival_time,
       service_route_availability.frequency,
       service_route_availability.ci_time_id,
       service_route_availability.notes,
       service_route_pricing.price_type,
       service_route_pricing.adult_price,
       service_route_pricing.child_price,
       service_route_pricing.infant_price,
       service_route_pricing.excluded_dates AS service_route_pricing_excluded_dates,
       service_route_pricing.sold_out_dates AS service_route_pricing_sold_out_dates,
       fare_type.name AS fare_type_name,
       start_station.name AS start_station_name,
       start_station.address AS start_station_address,
       end_station.name AS end_station_name,
       end_station.address AS end_station_address
FROM service_route
    INNER JOIN service ON service_route.service_id=service.id
    INNER JOIN product_type ON service.product_type_id=product_type.id
    INNER JOIN service_route_availability ON service_route_availability.service_route_id=service_route.id
    INNER JOIN service_route_pricing ON service_route_pricing.service_route_availability_id=service_route_availability.id
    INNER JOIN fare_type ON service_route_pricing.fare_type_id=fare_type.id
    INNER JOIN station start_station ON service_route.start_station_id=start_station.id
    INNER JOIN station end_station ON service_route.end_station_id=end_station.id
WHERE product_type.slug=sqlc.arg(product_type_slug)::TEXT
    AND service.deleted > now() AT TIME ZONE 'UTC'
    AND start_station.location_id=sqlc.arg(start_location_id)::UUID
    AND end_station.location_id=sqlc.arg(end_location_id)::UUID
    AND service_route_pricing.start_date<=sqlc.arg(search_start_date)::TIMESTAMP WITH TIME ZONE
    AND service_route_pricing.end_date>=sqlc.arg(search_end_date)::TIMESTAMP WITH TIME ZONE
    AND service_route_availability.start_date<=sqlc.arg(search_start_date)::TIMESTAMP WITH TIME ZONE
    AND service_route_availability.end_date>=sqlc.arg(search_end_date)::TIMESTAMP WITH TIME ZONE;


-- name: CreateServiceRouteAvailability :one
INSERT INTO service_route_availability (
    id,
    created,
    modified,
    service_route_id,
    departure_time,
    arrival_time,
    frequency,
    start_date,
    end_date,
    excluded_dates,
    sold_out_dates,
    ci_time_id,
    notes,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
         )
RETURNING *;

-- name: UpsertServiceRouteAvailability :one
INSERT INTO service_route_availability (
    id,
    created,
    modified,
    service_route_id,
    departure_time,
    arrival_time,
    frequency,
    start_date,
    end_date,
    excluded_dates,
    sold_out_dates,
    ci_time_id,
    notes,
    is_active
) VALUES (
             coalesce(nullif(@id, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
         )
ON CONFLICT (id) DO UPDATE
    SET
        modified = now() AT TIME ZONE 'UTC',
        service_route_id = $1,
        departure_time = $2,
        arrival_time = $3,
        frequency = $4,
        start_date = $5,
        end_date = $6,
        excluded_dates = $7,
        sold_out_dates = $8,
        ci_time_id = $9,
        notes = $10,
        is_active = $11
RETURNING *;

-- name: DeleteServiceRouteAvailability :exec
DELETE FROM service_route_availability
WHERE id = $1;

-- name: DeleteServiceRouteAvailabilityByServiceRouteId :exec
DELETE FROM service_route_availability
WHERE service_route_id = $1;

-- name: UpdateServiceRouteAvailability :one
UPDATE service_route_availability
SET
    modified = now() AT TIME ZONE 'UTC',
    service_route_id = $2,
    departure_time = $3,
    arrival_time = $4,
    frequency = $5,
    start_date = $6,
    end_date = $7,
    excluded_dates = $8,
    sold_out_dates = $9,
    ci_time_id = $10,
    notes = $11,
    is_active = $12
WHERE id = $1
RETURNING *;
