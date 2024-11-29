-- name: GetPackageById :one
SELECT * FROM package
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPackageBySlug :one
SELECT * FROM package
WHERE slug = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetPackageByCode :one
SELECT * FROM package
WHERE code = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountPackages :one
SELECT
    COUNT(id)
FROM package
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListPackagesAsc :many
SELECT * FROM package
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackagesDesc :many
SELECT * FROM package
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackage :one
INSERT INTO package (
    id,
    created,
    modified,
    deleted,
    title,
    slug,
    code,
    introduction,
    hero_image,
    hero_image_info,
    duration,
    required_deposit,
    departs,
    destinations,
    route,
    description,
    included_description,
    activities_description,
    journey_map,
    island_filter,
    keywords,
    is_group_journey,
    is_active
) VALUES (
     uuid_generate_v4(),
     now() AT TIME ZONE 'UTC',
     now() AT TIME ZONE 'UTC',
     'infinity'::timestamp AT TIME ZONE 'UTC',
     $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
 )
RETURNING *;

-- name: DeletePackage :exec
UPDATE package
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdatePackage :one
UPDATE package
SET
    modified = now() AT TIME ZONE 'UTC',
    title = $2,
    slug = $3,
    code = $4,
    introduction = $5,
    hero_image = $6,
    hero_image_info = $7,
    duration = $8,
    required_deposit = $9,
    departs = $10,
    destinations = $11,
    route = $12,
    description = $13,
    included_description = $14,
    activities_description = $15,
    journey_map = $16,
    island_filter = $17,
    keywords = $18,
    is_group_journey = $19,
    is_active = $20
WHERE id = $1
RETURNING *;
