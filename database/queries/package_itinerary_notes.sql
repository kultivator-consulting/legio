-- name: GetPackageItineraryNotesById :one
SELECT * FROM package_itinerary_notes
WHERE id = $1
LIMIT 1;

-- name: GetPackageItineraryNotesByPackageItineraryId :many
SELECT * FROM package_itinerary_notes
WHERE package_itinerary_id = $1;

-- name: CountPackageItineraryNotesByPackageItineraryId :one
SELECT
    COUNT(id)
FROM package_itinerary_notes
WHERE package_itinerary_id = $1;

-- name: ListPackageItineraryNotesByPackageItineraryIdAsc :many
SELECT * FROM package_itinerary_notes
WHERE package_itinerary_id = $1
ORDER BY sqlc.arg(sort_by)::text
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: ListPackageItineraryNotesByPackageItineraryIdDesc :many
SELECT * FROM package_itinerary_notes
WHERE package_itinerary_id = $1
ORDER BY sqlc.arg(sort_by)::text DESC
OFFSET (sqlc.arg(requested_page)::int - 1) * sqlc.arg(requested_page_size)::int
    FETCH NEXT sqlc.arg(requested_page_size) ROWS ONLY;

-- name: CreatePackageItineraryNotes :one
INSERT INTO package_itinerary_notes (
    id,
    created,
    modified,
    package_itinerary_id,
    note
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: UpdatePackageItineraryNotes :one
UPDATE package_itinerary_notes
SET
    modified = now() AT TIME ZONE 'UTC',
    package_itinerary_id = $2,
    note = $3
WHERE id = $1
RETURNING *;

-- name: DeletePackageItineraryNotes :exec
DELETE FROM package_itinerary_notes
WHERE id = $1;

-- name: DeletePackageItineraryNotesByPackageItineraryId :exec
DELETE FROM package_itinerary_notes
WHERE package_itinerary_id = $1;
