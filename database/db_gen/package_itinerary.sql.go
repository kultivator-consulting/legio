// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: package_itinerary.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countPackageItinerariesByPackageId = `-- name: CountPackageItinerariesByPackageId :one
SELECT
    COUNT(id)
FROM package_itinerary
WHERE package_id = $1
`

func (q *Queries) CountPackageItinerariesByPackageId(ctx context.Context, packageID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countPackageItinerariesByPackageId, packageID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPackageItinerary = `-- name: CreatePackageItinerary :one
INSERT INTO package_itinerary (
    id,
    created,
    modified,
    package_id,
    day,
    title,
    event_icon,
    station_id,
    latitude,
    longitude,
    description,
    supplier,
    supplier_code
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
         )
RETURNING id, created, modified, package_id, day, title, event_icon, station_id, latitude, longitude, description, supplier, supplier_code
`

type CreatePackageItineraryParams struct {
	PackageID    pgtype.UUID   `db:"package_id" json:"packageId"`
	Day          int32         `db:"day" json:"day"`
	Title        string        `db:"title" json:"title"`
	EventIcon    string        `db:"event_icon" json:"eventIcon"`
	StationID    pgtype.UUID   `db:"station_id" json:"stationId"`
	Latitude     pgtype.Float8 `db:"latitude" json:"latitude"`
	Longitude    pgtype.Float8 `db:"longitude" json:"longitude"`
	Description  string        `db:"description" json:"description"`
	Supplier     pgtype.Text   `db:"supplier" json:"supplier"`
	SupplierCode pgtype.Text   `db:"supplier_code" json:"supplierCode"`
}

func (q *Queries) CreatePackageItinerary(ctx context.Context, arg CreatePackageItineraryParams) (PackageItinerary, error) {
	row := q.db.QueryRow(ctx, createPackageItinerary,
		arg.PackageID,
		arg.Day,
		arg.Title,
		arg.EventIcon,
		arg.StationID,
		arg.Latitude,
		arg.Longitude,
		arg.Description,
		arg.Supplier,
		arg.SupplierCode,
	)
	var i PackageItinerary
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PackageID,
		&i.Day,
		&i.Title,
		&i.EventIcon,
		&i.StationID,
		&i.Latitude,
		&i.Longitude,
		&i.Description,
		&i.Supplier,
		&i.SupplierCode,
	)
	return i, err
}

const deletePackageItinerary = `-- name: DeletePackageItinerary :exec
DELETE FROM package_itinerary
WHERE id = $1
`

func (q *Queries) DeletePackageItinerary(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePackageItinerary, id)
	return err
}

const deletePackageItineraryByPackageId = `-- name: DeletePackageItineraryByPackageId :exec
DELETE FROM package_itinerary
WHERE package_id = $1
`

func (q *Queries) DeletePackageItineraryByPackageId(ctx context.Context, packageID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePackageItineraryByPackageId, packageID)
	return err
}

const getPackageItinerariesByPackageId = `-- name: GetPackageItinerariesByPackageId :many
SELECT id, created, modified, package_id, day, title, event_icon, station_id, latitude, longitude, description, supplier, supplier_code FROM package_itinerary
WHERE package_id = $1
`

func (q *Queries) GetPackageItinerariesByPackageId(ctx context.Context, packageID pgtype.UUID) ([]PackageItinerary, error) {
	rows, err := q.db.Query(ctx, getPackageItinerariesByPackageId, packageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PackageItinerary
	for rows.Next() {
		var i PackageItinerary
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.PackageID,
			&i.Day,
			&i.Title,
			&i.EventIcon,
			&i.StationID,
			&i.Latitude,
			&i.Longitude,
			&i.Description,
			&i.Supplier,
			&i.SupplierCode,
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

const getPackageItineraryById = `-- name: GetPackageItineraryById :one
SELECT id, created, modified, package_id, day, title, event_icon, station_id, latitude, longitude, description, supplier, supplier_code FROM package_itinerary
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetPackageItineraryById(ctx context.Context, id pgtype.UUID) (PackageItinerary, error) {
	row := q.db.QueryRow(ctx, getPackageItineraryById, id)
	var i PackageItinerary
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PackageID,
		&i.Day,
		&i.Title,
		&i.EventIcon,
		&i.StationID,
		&i.Latitude,
		&i.Longitude,
		&i.Description,
		&i.Supplier,
		&i.SupplierCode,
	)
	return i, err
}

const listPackageItinerariesByPackageIdAsc = `-- name: ListPackageItinerariesByPackageIdAsc :many
SELECT id, created, modified, package_id, day, title, event_icon, station_id, latitude, longitude, description, supplier, supplier_code FROM package_itinerary
WHERE package_id = $1
ORDER BY $2::text
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPackageItinerariesByPackageIdAscParams struct {
	PackageID         pgtype.UUID `db:"package_id" json:"packageId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPackageItinerariesByPackageIdAsc(ctx context.Context, arg ListPackageItinerariesByPackageIdAscParams) ([]PackageItinerary, error) {
	rows, err := q.db.Query(ctx, listPackageItinerariesByPackageIdAsc,
		arg.PackageID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PackageItinerary
	for rows.Next() {
		var i PackageItinerary
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.PackageID,
			&i.Day,
			&i.Title,
			&i.EventIcon,
			&i.StationID,
			&i.Latitude,
			&i.Longitude,
			&i.Description,
			&i.Supplier,
			&i.SupplierCode,
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

const listPackageItinerariesByPackageIdDesc = `-- name: ListPackageItinerariesByPackageIdDesc :many
SELECT id, created, modified, package_id, day, title, event_icon, station_id, latitude, longitude, description, supplier, supplier_code FROM package_itinerary
WHERE package_id = $1
ORDER BY $2::text DESC
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPackageItinerariesByPackageIdDescParams struct {
	PackageID         pgtype.UUID `db:"package_id" json:"packageId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPackageItinerariesByPackageIdDesc(ctx context.Context, arg ListPackageItinerariesByPackageIdDescParams) ([]PackageItinerary, error) {
	rows, err := q.db.Query(ctx, listPackageItinerariesByPackageIdDesc,
		arg.PackageID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PackageItinerary
	for rows.Next() {
		var i PackageItinerary
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.PackageID,
			&i.Day,
			&i.Title,
			&i.EventIcon,
			&i.StationID,
			&i.Latitude,
			&i.Longitude,
			&i.Description,
			&i.Supplier,
			&i.SupplierCode,
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

const updatePackageItinerary = `-- name: UpdatePackageItinerary :one
UPDATE package_itinerary
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    day = $3,
    title = $4,
    event_icon = $5,
    station_id = $6,
    latitude = $7,
    longitude = $8,
    description = $9,
    supplier = $10,
    supplier_code = $11
WHERE id = $1
RETURNING id, created, modified, package_id, day, title, event_icon, station_id, latitude, longitude, description, supplier, supplier_code
`

type UpdatePackageItineraryParams struct {
	ID           pgtype.UUID   `db:"id" json:"id"`
	PackageID    pgtype.UUID   `db:"package_id" json:"packageId"`
	Day          int32         `db:"day" json:"day"`
	Title        string        `db:"title" json:"title"`
	EventIcon    string        `db:"event_icon" json:"eventIcon"`
	StationID    pgtype.UUID   `db:"station_id" json:"stationId"`
	Latitude     pgtype.Float8 `db:"latitude" json:"latitude"`
	Longitude    pgtype.Float8 `db:"longitude" json:"longitude"`
	Description  string        `db:"description" json:"description"`
	Supplier     pgtype.Text   `db:"supplier" json:"supplier"`
	SupplierCode pgtype.Text   `db:"supplier_code" json:"supplierCode"`
}

func (q *Queries) UpdatePackageItinerary(ctx context.Context, arg UpdatePackageItineraryParams) (PackageItinerary, error) {
	row := q.db.QueryRow(ctx, updatePackageItinerary,
		arg.ID,
		arg.PackageID,
		arg.Day,
		arg.Title,
		arg.EventIcon,
		arg.StationID,
		arg.Latitude,
		arg.Longitude,
		arg.Description,
		arg.Supplier,
		arg.SupplierCode,
	)
	var i PackageItinerary
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PackageID,
		&i.Day,
		&i.Title,
		&i.EventIcon,
		&i.StationID,
		&i.Latitude,
		&i.Longitude,
		&i.Description,
		&i.Supplier,
		&i.SupplierCode,
	)
	return i, err
}
