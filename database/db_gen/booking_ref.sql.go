// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: booking_ref.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countBookingRefs = `-- name: CountBookingRefs :one
SELECT
    COUNT(id)
FROM booking_ref
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountBookingRefs(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countBookingRefs)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createBookingRef = `-- name: CreateBookingRef :one
INSERT INTO booking_ref (
    id,
    created,
    modified,
    deleted,
    prefix,
    initial,
    length,
    sequence
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4
         )
RETURNING id, created, modified, deleted, prefix, initial, length, sequence
`

type CreateBookingRefParams struct {
	Prefix   string `db:"prefix" json:"prefix"`
	Initial  int64  `db:"initial" json:"initial"`
	Length   int32  `db:"length" json:"length"`
	Sequence int64  `db:"sequence" json:"sequence"`
}

func (q *Queries) CreateBookingRef(ctx context.Context, arg CreateBookingRefParams) (BookingRef, error) {
	row := q.db.QueryRow(ctx, createBookingRef,
		arg.Prefix,
		arg.Initial,
		arg.Length,
		arg.Sequence,
	)
	var i BookingRef
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Prefix,
		&i.Initial,
		&i.Length,
		&i.Sequence,
	)
	return i, err
}

const deleteBookingRef = `-- name: DeleteBookingRef :exec
UPDATE booking_ref
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
`

func (q *Queries) DeleteBookingRef(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteBookingRef, id)
	return err
}

const getBookingRefById = `-- name: GetBookingRefById :one
SELECT id, created, modified, deleted, prefix, initial, length, sequence FROM booking_ref
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetBookingRefById(ctx context.Context, id pgtype.UUID) (BookingRef, error) {
	row := q.db.QueryRow(ctx, getBookingRefById, id)
	var i BookingRef
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Prefix,
		&i.Initial,
		&i.Length,
		&i.Sequence,
	)
	return i, err
}

const getBookingRefByPrefix = `-- name: GetBookingRefByPrefix :one
SELECT id, created, modified, deleted, prefix, initial, length, sequence FROM booking_ref
WHERE prefix = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetBookingRefByPrefix(ctx context.Context, prefix string) (BookingRef, error) {
	row := q.db.QueryRow(ctx, getBookingRefByPrefix, prefix)
	var i BookingRef
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Prefix,
		&i.Initial,
		&i.Length,
		&i.Sequence,
	)
	return i, err
}

const incrementBookingRefSequenceById = `-- name: IncrementBookingRefSequenceById :one
UPDATE booking_ref
SET
    sequence = sequence + 1
WHERE id = $1
RETURNING id, created, modified, deleted, prefix, initial, length, sequence
`

func (q *Queries) IncrementBookingRefSequenceById(ctx context.Context, id pgtype.UUID) (BookingRef, error) {
	row := q.db.QueryRow(ctx, incrementBookingRefSequenceById, id)
	var i BookingRef
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Prefix,
		&i.Initial,
		&i.Length,
		&i.Sequence,
	)
	return i, err
}

const incrementBookingRefSequenceByPrefix = `-- name: IncrementBookingRefSequenceByPrefix :one
UPDATE booking_ref
SET
    sequence = sequence + 1
WHERE prefix = $1
RETURNING id, created, modified, deleted, prefix, initial, length, sequence
`

func (q *Queries) IncrementBookingRefSequenceByPrefix(ctx context.Context, prefix string) (BookingRef, error) {
	row := q.db.QueryRow(ctx, incrementBookingRefSequenceByPrefix, prefix)
	var i BookingRef
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Prefix,
		&i.Initial,
		&i.Length,
		&i.Sequence,
	)
	return i, err
}

const listBookingRefsAsc = `-- name: ListBookingRefsAsc :many
SELECT id, created, modified, deleted, prefix, initial, length, sequence FROM booking_ref
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListBookingRefsAscParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListBookingRefsAsc(ctx context.Context, arg ListBookingRefsAscParams) ([]BookingRef, error) {
	rows, err := q.db.Query(ctx, listBookingRefsAsc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BookingRef
	for rows.Next() {
		var i BookingRef
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.Prefix,
			&i.Initial,
			&i.Length,
			&i.Sequence,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listBookingRefsDesc = `-- name: ListBookingRefsDesc :many
SELECT id, created, modified, deleted, prefix, initial, length, sequence FROM booking_ref
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY $1::text DESC
OFFSET ($2::int - 1) * $3::int
    FETCH NEXT $3 ROWS ONLY
`

type ListBookingRefsDescParams struct {
	SortBy            string `db:"sort_by" json:"sortBy"`
	RequestedPage     int32  `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32  `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListBookingRefsDesc(ctx context.Context, arg ListBookingRefsDescParams) ([]BookingRef, error) {
	rows, err := q.db.Query(ctx, listBookingRefsDesc, arg.SortBy, arg.RequestedPage, arg.RequestedPageSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []BookingRef
	for rows.Next() {
		var i BookingRef
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.Prefix,
			&i.Initial,
			&i.Length,
			&i.Sequence,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBookingRef = `-- name: UpdateBookingRef :one
UPDATE booking_ref
SET
    modified = now() AT TIME ZONE 'UTC',
    prefix = $2,
    initial = $3,
    length = $4,
    sequence = $5
WHERE id = $1
RETURNING id, created, modified, deleted, prefix, initial, length, sequence
`

type UpdateBookingRefParams struct {
	ID       pgtype.UUID `db:"id" json:"id"`
	Prefix   string      `db:"prefix" json:"prefix"`
	Initial  int64       `db:"initial" json:"initial"`
	Length   int32       `db:"length" json:"length"`
	Sequence int64       `db:"sequence" json:"sequence"`
}

func (q *Queries) UpdateBookingRef(ctx context.Context, arg UpdateBookingRefParams) (BookingRef, error) {
	row := q.db.QueryRow(ctx, updateBookingRef,
		arg.ID,
		arg.Prefix,
		arg.Initial,
		arg.Length,
		arg.Sequence,
	)
	var i BookingRef
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.Prefix,
		&i.Initial,
		&i.Length,
		&i.Sequence,
	)
	return i, err
}
