// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: pass_availability_pricing.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countPassAvailabilityPricing = `-- name: CountPassAvailabilityPricing :one
SELECT
    COUNT(id)
FROM pass_availability_pricing
WHERE pass_availability_id = $1
`

func (q *Queries) CountPassAvailabilityPricing(ctx context.Context, passAvailabilityID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countPassAvailabilityPricing, passAvailabilityID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPassAvailabilityPricing = `-- name: CreatePassAvailabilityPricing :one
INSERT INTO pass_availability_pricing (
    id,
    created,
    modified,
    pass_availability_id,
    fare_type_id,
    adult_price,
    child_price,
    infant_price,
    start_date,
    end_date,
    excluded_dates
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8
         )
RETURNING id, created, modified, pass_availability_id, fare_type_id, adult_price, child_price, infant_price, start_date, end_date, excluded_dates
`

type CreatePassAvailabilityPricingParams struct {
	PassAvailabilityID pgtype.UUID          `db:"pass_availability_id" json:"passAvailabilityId"`
	FareTypeID         pgtype.UUID          `db:"fare_type_id" json:"fareTypeId"`
	AdultPrice         pgtype.Numeric       `db:"adult_price" json:"adultPrice"`
	ChildPrice         pgtype.Numeric       `db:"child_price" json:"childPrice"`
	InfantPrice        pgtype.Numeric       `db:"infant_price" json:"infantPrice"`
	StartDate          pgtype.Timestamptz   `db:"start_date" json:"startDate"`
	EndDate            pgtype.Timestamptz   `db:"end_date" json:"endDate"`
	ExcludedDates      []pgtype.Timestamptz `db:"excluded_dates" json:"excludedDates"`
}

func (q *Queries) CreatePassAvailabilityPricing(ctx context.Context, arg CreatePassAvailabilityPricingParams) (PassAvailabilityPricing, error) {
	row := q.db.QueryRow(ctx, createPassAvailabilityPricing,
		arg.PassAvailabilityID,
		arg.FareTypeID,
		arg.AdultPrice,
		arg.ChildPrice,
		arg.InfantPrice,
		arg.StartDate,
		arg.EndDate,
		arg.ExcludedDates,
	)
	var i PassAvailabilityPricing
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PassAvailabilityID,
		&i.FareTypeID,
		&i.AdultPrice,
		&i.ChildPrice,
		&i.InfantPrice,
		&i.StartDate,
		&i.EndDate,
		&i.ExcludedDates,
	)
	return i, err
}

const deletePassAvailabilityPricing = `-- name: DeletePassAvailabilityPricing :exec
DELETE FROM pass_availability_pricing
WHERE id = $1
`

func (q *Queries) DeletePassAvailabilityPricing(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePassAvailabilityPricing, id)
	return err
}

const getPassAvailabilityPricingById = `-- name: GetPassAvailabilityPricingById :one
SELECT id, created, modified, pass_availability_id, fare_type_id, adult_price, child_price, infant_price, start_date, end_date, excluded_dates FROM pass_availability_pricing
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetPassAvailabilityPricingById(ctx context.Context, id pgtype.UUID) (PassAvailabilityPricing, error) {
	row := q.db.QueryRow(ctx, getPassAvailabilityPricingById, id)
	var i PassAvailabilityPricing
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PassAvailabilityID,
		&i.FareTypeID,
		&i.AdultPrice,
		&i.ChildPrice,
		&i.InfantPrice,
		&i.StartDate,
		&i.EndDate,
		&i.ExcludedDates,
	)
	return i, err
}

const listPassAvailabilityPricingByPassAvailabilityId = `-- name: ListPassAvailabilityPricingByPassAvailabilityId :many
SELECT id, created, modified, pass_availability_id, fare_type_id, adult_price, child_price, infant_price, start_date, end_date, excluded_dates FROM pass_availability_pricing
WHERE pass_availability_id = $1
`

func (q *Queries) ListPassAvailabilityPricingByPassAvailabilityId(ctx context.Context, passAvailabilityID pgtype.UUID) ([]PassAvailabilityPricing, error) {
	rows, err := q.db.Query(ctx, listPassAvailabilityPricingByPassAvailabilityId, passAvailabilityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PassAvailabilityPricing
	for rows.Next() {
		var i PassAvailabilityPricing
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.PassAvailabilityID,
			&i.FareTypeID,
			&i.AdultPrice,
			&i.ChildPrice,
			&i.InfantPrice,
			&i.StartDate,
			&i.EndDate,
			&i.ExcludedDates,
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

const updatePassAvailabilityPricing = `-- name: UpdatePassAvailabilityPricing :one
UPDATE pass_availability_pricing
SET
    modified = now() AT TIME ZONE 'UTC',
    pass_availability_id = $2,
    fare_type_id = $3,
    adult_price = $4,
    child_price = $5,
    infant_price = $6,
    start_date = $7,
    end_date = $8,
    excluded_dates = $9
WHERE id = $1
RETURNING id, created, modified, pass_availability_id, fare_type_id, adult_price, child_price, infant_price, start_date, end_date, excluded_dates
`

type UpdatePassAvailabilityPricingParams struct {
	ID                 pgtype.UUID          `db:"id" json:"id"`
	PassAvailabilityID pgtype.UUID          `db:"pass_availability_id" json:"passAvailabilityId"`
	FareTypeID         pgtype.UUID          `db:"fare_type_id" json:"fareTypeId"`
	AdultPrice         pgtype.Numeric       `db:"adult_price" json:"adultPrice"`
	ChildPrice         pgtype.Numeric       `db:"child_price" json:"childPrice"`
	InfantPrice        pgtype.Numeric       `db:"infant_price" json:"infantPrice"`
	StartDate          pgtype.Timestamptz   `db:"start_date" json:"startDate"`
	EndDate            pgtype.Timestamptz   `db:"end_date" json:"endDate"`
	ExcludedDates      []pgtype.Timestamptz `db:"excluded_dates" json:"excludedDates"`
}

func (q *Queries) UpdatePassAvailabilityPricing(ctx context.Context, arg UpdatePassAvailabilityPricingParams) (PassAvailabilityPricing, error) {
	row := q.db.QueryRow(ctx, updatePassAvailabilityPricing,
		arg.ID,
		arg.PassAvailabilityID,
		arg.FareTypeID,
		arg.AdultPrice,
		arg.ChildPrice,
		arg.InfantPrice,
		arg.StartDate,
		arg.EndDate,
		arg.ExcludedDates,
	)
	var i PassAvailabilityPricing
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PassAvailabilityID,
		&i.FareTypeID,
		&i.AdultPrice,
		&i.ChildPrice,
		&i.InfantPrice,
		&i.StartDate,
		&i.EndDate,
		&i.ExcludedDates,
	)
	return i, err
}

const upsertPassAvailabilityPricing = `-- name: UpsertPassAvailabilityPricing :one
INSERT INTO pass_availability_pricing (
    id,
    created,
    modified,
    pass_availability_id,
    fare_type_id,
    adult_price,
    child_price,
    infant_price,
    start_date,
    end_date,
    excluded_dates
) VALUES (
             coalesce(nullif($9, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
                $1, $2, $3, $4, $5, $6, $7, $8
         )
ON CONFLICT (id) DO UPDATE
    SET
        modified = now() AT TIME ZONE 'UTC',
        pass_availability_id = $1,
        fare_type_id = $2,
        adult_price = $3,
        child_price = $4,
        infant_price = $5,
        start_date = $6,
        end_date = $7,
        excluded_dates = $8
RETURNING id, created, modified, pass_availability_id, fare_type_id, adult_price, child_price, infant_price, start_date, end_date, excluded_dates
`

type UpsertPassAvailabilityPricingParams struct {
	PassAvailabilityID pgtype.UUID          `db:"pass_availability_id" json:"passAvailabilityId"`
	FareTypeID         pgtype.UUID          `db:"fare_type_id" json:"fareTypeId"`
	AdultPrice         pgtype.Numeric       `db:"adult_price" json:"adultPrice"`
	ChildPrice         pgtype.Numeric       `db:"child_price" json:"childPrice"`
	InfantPrice        pgtype.Numeric       `db:"infant_price" json:"infantPrice"`
	StartDate          pgtype.Timestamptz   `db:"start_date" json:"startDate"`
	EndDate            pgtype.Timestamptz   `db:"end_date" json:"endDate"`
	ExcludedDates      []pgtype.Timestamptz `db:"excluded_dates" json:"excludedDates"`
	ID                 interface{}          `db:"id" json:"id"`
}

func (q *Queries) UpsertPassAvailabilityPricing(ctx context.Context, arg UpsertPassAvailabilityPricingParams) (PassAvailabilityPricing, error) {
	row := q.db.QueryRow(ctx, upsertPassAvailabilityPricing,
		arg.PassAvailabilityID,
		arg.FareTypeID,
		arg.AdultPrice,
		arg.ChildPrice,
		arg.InfantPrice,
		arg.StartDate,
		arg.EndDate,
		arg.ExcludedDates,
		arg.ID,
	)
	var i PassAvailabilityPricing
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PassAvailabilityID,
		&i.FareTypeID,
		&i.AdultPrice,
		&i.ChildPrice,
		&i.InfantPrice,
		&i.StartDate,
		&i.EndDate,
		&i.ExcludedDates,
	)
	return i, err
}
