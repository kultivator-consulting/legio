-- name: GetPackageIncludedActivitiesById :one
SELECT * FROM package_included_activities
WHERE id = $1
LIMIT 1;

-- name: GetPackageIncludedActivitiesByPackageId :many
SELECT * FROM package_included_activities
WHERE package_id = $1;

-- name: CountPackageIncludedActivitiesByPackageId :one
SELECT
    COUNT(id)
FROM package_included_activities
WHERE package_id = $1;

-- name: ListPackageIncludedActivitiesByPackageIdAsc :many
SELECT * FROM package_included_activities
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageIncludedActivitiesByPackageIdDesc :many
SELECT * FROM package_included_activities
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackageIncludedActivities :one
INSERT INTO package_included_activities (
    id,
    created,
    modified,
    package_id,
    ordering,
    title,
    description,
    hero_image,
    hero_image_info,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7
         )
RETURNING *;

-- name: UpdatePackageIncludedActivities :one
UPDATE package_included_activities
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    ordering = $3,
    title = $4,
    description = $5,
    hero_image = $6,
    hero_image_info = $7,
    is_active = $8
WHERE id = $1
RETURNING *;

-- name: DeletePackageIncludedActivities :exec
DELETE FROM package_included_activities
WHERE id = $1;

-- name: DeletePackageIncludedActivitiesByPackageId :exec
DELETE FROM package_included_activities
WHERE package_id = $1;
