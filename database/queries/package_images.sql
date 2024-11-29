-- name: GetPackageImageById :one
SELECT * FROM package_images
WHERE id = $1
LIMIT 1;

-- name: GetPackageImagesByPackageId :many
SELECT * FROM package_images
WHERE package_id = $1;

-- name: CountPackageImagesByPackageId :one
SELECT
    COUNT(id)
FROM package_images
WHERE package_id = $1;

-- name: ListPackageImagesByPackageIdAsc :many
SELECT * FROM package_images
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageImagesByPackageIdDesc :many
SELECT * FROM package_images
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackageImage :one
INSERT INTO package_images (
    id,
    created,
    modified,
    package_id,
    ordering,
    image,
    image_info
) VALUES (
     uuid_generate_v4(),
     now() AT TIME ZONE 'UTC',
     now() AT TIME ZONE 'UTC',
     $1, $2, $3, $4
 )
RETURNING *;

-- name: UpdatePackageImage :one
UPDATE package_images
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    ordering = $3,
    image = $4,
    image_info = $5
WHERE id = $1
RETURNING *;

-- name: DeletePackageImage :exec
DELETE FROM package_images
WHERE id = $1;

-- name: DeletePackageImagesByPackageId :exec
DELETE FROM package_images
WHERE package_id = $1;
