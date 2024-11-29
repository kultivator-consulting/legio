-- name: GetServiceById :one
SELECT * FROM service
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetServiceByName :one
SELECT * FROM service
WHERE name = $1
  AND deleted > now() AT TIME ZONE 'UTC'
  AND is_active = TRUE
LIMIT 1;

-- name: SearchServicesByOperatorIdCount :one
SELECT
    COUNT(id)
FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND operator_id = $1
  AND name ILIKE '%' || sqlc.arg(search) || '%'
   OR sub_service_name ILIKE '%' || sqlc.arg(search) || '%'
   OR description ILIKE '%' || sqlc.arg(search) || '%'
   OR tour_name ILIKE '%' || sqlc.arg(search) || '%'
   OR service_code ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_title ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_content ILIKE '%' || sqlc.arg(search) || '%'
   OR excursion_text ILIKE '%' || sqlc.arg(search) || '%';

-- name: SearchServicesByOperatorIdAsc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND operator_id = $1
  AND name ILIKE '%' || sqlc.arg(search) || '%'
   OR sub_service_name ILIKE '%' || sqlc.arg(search) || '%'
   OR description ILIKE '%' || sqlc.arg(search) || '%'
   OR tour_name ILIKE '%' || sqlc.arg(search) || '%'
   OR service_code ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_title ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_content ILIKE '%' || sqlc.arg(search) || '%'
   OR excursion_text ILIKE '%' || sqlc.arg(search) || '%'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: SearchServicesByOperatorIdDesc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND operator_id = $1
  AND name ILIKE '%' || sqlc.arg(search) || '%'
   OR sub_service_name ILIKE '%' || sqlc.arg(search) || '%'
   OR description ILIKE '%' || sqlc.arg(search) || '%'
   OR tour_name ILIKE '%' || sqlc.arg(search) || '%'
   OR service_code ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_title ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_content ILIKE '%' || sqlc.arg(search) || '%'
   OR excursion_text ILIKE '%' || sqlc.arg(search) || '%'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: SearchServicesByProductTypeIdCount :one
SELECT
    COUNT(id)
FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND product_type_id = $1
  AND name ILIKE '%' || sqlc.arg(search) || '%'
   OR sub_service_name ILIKE '%' || sqlc.arg(search) || '%'
   OR description ILIKE '%' || sqlc.arg(search) || '%'
   OR tour_name ILIKE '%' || sqlc.arg(search) || '%'
   OR service_code ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_title ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_content ILIKE '%' || sqlc.arg(search) || '%'
   OR excursion_text ILIKE '%' || sqlc.arg(search) || '%';

-- name: SearchServicesByProductTypeAsc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND product_type_id = $1
  AND name ILIKE '%' || sqlc.arg(search) || '%'
   OR sub_service_name ILIKE '%' || sqlc.arg(search) || '%'
   OR description ILIKE '%' || sqlc.arg(search) || '%'
   OR tour_name ILIKE '%' || sqlc.arg(search) || '%'
   OR service_code ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_title ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_content ILIKE '%' || sqlc.arg(search) || '%'
   OR excursion_text ILIKE '%' || sqlc.arg(search) || '%'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: SearchServicesByProductTypeDesc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND product_type_id = $1
  AND name ILIKE '%' || sqlc.arg(search) || '%'
    OR sub_service_name ILIKE '%' || sqlc.arg(search) || '%'
    OR description ILIKE '%' || sqlc.arg(search) || '%'
    OR tour_name ILIKE '%' || sqlc.arg(search) || '%'
    OR service_code ILIKE '%' || sqlc.arg(search) || '%'
    OR meta_title ILIKE '%' || sqlc.arg(search) || '%'
    OR meta_content ILIKE '%' || sqlc.arg(search) || '%'
    OR excursion_text ILIKE '%' || sqlc.arg(search) || '%'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: SearchServicesCount :one
SELECT
    COUNT(id)
FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND name ILIKE '%' || sqlc.arg(search) || '%'
   OR sub_service_name ILIKE '%' || sqlc.arg(search) || '%'
   OR description ILIKE '%' || sqlc.arg(search) || '%'
   OR tour_name ILIKE '%' || sqlc.arg(search) || '%'
   OR service_code ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_title ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_content ILIKE '%' || sqlc.arg(search) || '%'
   OR excursion_text ILIKE '%' || sqlc.arg(search) || '%';

-- name: SearchServicesAsc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND name ILIKE '%' || sqlc.arg(search) || '%'
   OR sub_service_name ILIKE '%' || sqlc.arg(search) || '%'
   OR description ILIKE '%' || sqlc.arg(search) || '%'
   OR tour_name ILIKE '%' || sqlc.arg(search) || '%'
   OR service_code ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_title ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_content ILIKE '%' || sqlc.arg(search) || '%'
   OR excursion_text ILIKE '%' || sqlc.arg(search) || '%'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: SearchServicesDesc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
  AND name ILIKE '%' || sqlc.arg(search) || '%'
   OR sub_service_name ILIKE '%' || sqlc.arg(search) || '%'
   OR description ILIKE '%' || sqlc.arg(search) || '%'
   OR tour_name ILIKE '%' || sqlc.arg(search) || '%'
   OR service_code ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_title ILIKE '%' || sqlc.arg(search) || '%'
   OR meta_content ILIKE '%' || sqlc.arg(search) || '%'
   OR excursion_text ILIKE '%' || sqlc.arg(search) || '%'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CountServicesByOperatorId :one
SELECT
    COUNT(id)
FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND operator_id = $1;

-- name: ListServicesByOperatorIdAsc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND operator_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListServicesByOperatorIdDesc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND operator_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CountServicesByProductTypeId :one
SELECT
    COUNT(id)
FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND product_type_id = $1;

-- name: ListServicesByProductTypeAsc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND product_type_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListServicesByProductTypeDesc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
    AND product_type_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CountServices :one
SELECT
    COUNT(id)
FROM service
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListServicesAsc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListServicesDesc :many
SELECT * FROM service
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreateService :one
INSERT INTO service (
    id,
    created,
    modified,
    deleted,
    operator_id,
    product_type_id,
    name,
    sub_service_name,
    start_location_id,
    end_location_id,
    start_station_id,
    end_station_id,
    logo,
    banner_image,
    excluded_dates,
    sold_out_dates,
    is_popular,
    with_pass,
    is_tour,
    tour_name,
    description,
    day_excursion,
    multi_service,
    excursion_text,
    service_code,
    meta_title,
    meta_content,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6,
             $7, $8, $9, $10, $11, $12,
             $13, $14, $15, $16, $17, $18, $19,
             $20, $21, $22, $23, $24
         )
RETURNING *;

-- name: DeleteService :exec
UPDATE service
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateService :one
UPDATE service
SET
    modified = now() AT TIME ZONE 'UTC',
    operator_id = $2,
    product_type_id = $3,
    name = $4,
    sub_service_name = $5,
    start_location_id = $6,
    end_location_id = $7,
    start_station_id = $8,
    end_station_id = $9,
    logo = $10,
    banner_image = $11,
    excluded_dates = $12,
    sold_out_dates = $13,
    is_popular = $14,
    with_pass = $15,
    is_tour = $16,
    tour_name = $17,
    description = $18,
    day_excursion = $19,
    multi_service = $20,
    excursion_text = $21,
    service_code = $22,
    meta_title = $23,
    meta_content = $24,
    is_active = $25
WHERE id = $1
RETURNING *;
