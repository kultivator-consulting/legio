-- name: GetStationProductTypeById :one
SELECT * FROM station_product_type
WHERE id = $1
LIMIT 1;

-- name: GetStationProductTypeByStationId :many
SELECT * FROM station_product_type
WHERE station_id = $1;

-- name: CountStationProductTypesByStationId :one
SELECT
    COUNT(id)
FROM station_product_type
WHERE station_id = $1;

-- name: ListStationProductTypesByStationIdAsc :many
SELECT * FROM station_product_type
WHERE station_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListStationProductTypesByStationIdDesc :many
SELECT * FROM station_product_type
WHERE station_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateStationProductType :one
INSERT INTO station_product_type (
    id,
    created,
    modified,
    station_id,
    product_type_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: UpdateStationProductType :one
UPDATE station_product_type
SET
    modified = now() AT TIME ZONE 'UTC',
    station_id = $2,
    product_type_id = $3
WHERE id = $1
RETURNING *;

-- name: DeleteStationProductType :exec
DELETE FROM station_product_type
WHERE id = $1;

-- name: DeleteStationProductTypesByStationId :exec
DELETE FROM station_product_type
WHERE station_id = $1;
