-- name: GetExtraOperatorById :one
SELECT * FROM extra_operator
WHERE id = $1
LIMIT 1;

-- name: CountExtraOperatorsByExtraId :one
SELECT
    COUNT(id)
FROM extra_operator
WHERE extra_id = $1;

-- name: ListExtraOperatorsByExtraId :many
SELECT * FROM extra_operator
WHERE extra_id = $1;

-- name: CountExtraOperatorsByOperatorId :one
SELECT
    COUNT(id)
FROM extra_operator
WHERE operator_id = $1;

-- name: ListExtraOperatorsByOperatorId :many
SELECT * FROM extra_operator
WHERE operator_id = $1;

-- name: CreateExtraOperator :one
INSERT INTO extra_operator (
    id,
    created,
    modified,
    extra_id,
    operator_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING *;

-- name: DeleteExtraOperator :exec
DELETE FROM extra_operator
WHERE id = $1;

-- name: DeleteExtraOperatorByExtraId :exec
DELETE FROM extra_operator
WHERE extra_id = $1;

-- name: UpdateExtraOperator :one
UPDATE extra_operator
SET
    modified = now() AT TIME ZONE 'UTC',
    extra_id = $2,
    operator_id = $3
WHERE id = $1
RETURNING *;
