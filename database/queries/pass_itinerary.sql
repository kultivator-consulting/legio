-- name: GetPassItineraryById :one
SELECT * FROM pass_itinerary
WHERE id = $1
LIMIT 1;

-- name: CountPassItineraries :one
SELECT
    COUNT(id)
FROM pass_itinerary
WHERE pass_id = $1;

-- name: ListPassItinerariesByPassId :many
SELECT * FROM pass_itinerary
WHERE pass_id = $1;

-- name: CreatePassItinerary :one
INSERT INTO pass_itinerary (
    id,
    created,
    modified,
    pass_id,
    title,
    description,
    terms
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING *;

-- name: UpsertPassItinerary :one
INSERT INTO pass_itinerary (
    id,
    created,
    modified,
    pass_id,
    title,
    description,
    terms
) VALUES (
             coalesce(nullif(@id, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
                $1, $2, $3, $4
         )
ON CONFLICT (id) DO UPDATE
    SET
        modified = now() AT TIME ZONE 'UTC',
        pass_id = $1,
        title = $2,
        description = $3,
        terms = $4
RETURNING *;

-- name: DeletePassItinerary :exec
DELETE FROM pass_itinerary
WHERE id = $1;

-- name: UpdatePassItinerary :one
UPDATE pass_itinerary
SET
    modified = now() AT TIME ZONE 'UTC',
    pass_id = $2,
    title = $3,
    description = $4,
    terms = $5
WHERE id = $1
RETURNING *;
