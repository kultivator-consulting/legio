-- name: GetPackageIncludedItemsById :one
SELECT * FROM package_included_items
WHERE id = $1
LIMIT 1;

-- name: GetPackageIncludedItemsByPackageId :many
SELECT * FROM package_included_items
WHERE package_id = $1;

-- name: CountPackageIncludedItemsByPackageId :one
SELECT
    COUNT(id)
FROM package_included_items
WHERE package_id = $1;

-- name: ListPackageIncludedItemsByPackageIdAsc :many
SELECT * FROM package_included_items
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageIncludedItemsByPackageIdDesc :many
SELECT * FROM package_included_items
WHERE package_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackageIncludedItems :one
INSERT INTO package_included_items (
    id,
    created,
    modified,
    package_id,
    ordering,
    type,
    description,
    item_icon,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6
         )
RETURNING *;

-- name: UpdatePackageIncludedItems :one
UPDATE package_included_items
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    ordering = $3,
    type = $4,
    description = $5,
    item_icon = $6,
    is_active = $7
WHERE id = $1
RETURNING *;

-- name: DeletePackageIncludedItems :exec
DELETE FROM package_included_items
WHERE id = $1;

-- name: DeletePackageIncludedItemsByPackageId :exec
DELETE FROM package_included_items
WHERE package_id = $1;
