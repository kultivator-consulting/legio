// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: service_route.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countServiceRoutesByServiceId = `-- name: CountServiceRoutesByServiceId :one
SELECT
    COUNT(id)
FROM service_route
WHERE service_id = $1
`

func (q *Queries) CountServiceRoutesByServiceId(ctx context.Context, serviceID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countServiceRoutesByServiceId, serviceID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createServiceRoute = `-- name: CreateServiceRoute :one
INSERT INTO service_route (
    id,
    created,
    modified,
    service_id,
    product_id,
    start_station_id,
    end_station_id,
    is_main_route,
    is_popular,
    ic_route_id,
    short_name,
    description,
    url,
    color,
    text_color
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
         )
RETURNING id, created, modified, service_id, product_id, start_station_id, end_station_id, is_main_route, is_popular, ic_route_id, short_name, description, url, color, text_color, is_active
`

type CreateServiceRouteParams struct {
	ServiceID      pgtype.UUID `db:"service_id" json:"serviceId"`
	ProductID      pgtype.UUID `db:"product_id" json:"productId"`
	StartStationID pgtype.UUID `db:"start_station_id" json:"startStationId"`
	EndStationID   pgtype.UUID `db:"end_station_id" json:"endStationId"`
	IsMainRoute    bool        `db:"is_main_route" json:"isMainRoute"`
	IsPopular      bool        `db:"is_popular" json:"isPopular"`
	IcRouteID      pgtype.Text `db:"ic_route_id" json:"icRouteId"`
	ShortName      pgtype.Text `db:"short_name" json:"shortName"`
	Description    pgtype.Text `db:"description" json:"description"`
	Url            pgtype.Text `db:"url" json:"url"`
	Color          pgtype.Text `db:"color" json:"color"`
	TextColor      pgtype.Text `db:"text_color" json:"textColor"`
}

func (q *Queries) CreateServiceRoute(ctx context.Context, arg CreateServiceRouteParams) (ServiceRoute, error) {
	row := q.db.QueryRow(ctx, createServiceRoute,
		arg.ServiceID,
		arg.ProductID,
		arg.StartStationID,
		arg.EndStationID,
		arg.IsMainRoute,
		arg.IsPopular,
		arg.IcRouteID,
		arg.ShortName,
		arg.Description,
		arg.Url,
		arg.Color,
		arg.TextColor,
	)
	var i ServiceRoute
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ServiceID,
		&i.ProductID,
		&i.StartStationID,
		&i.EndStationID,
		&i.IsMainRoute,
		&i.IsPopular,
		&i.IcRouteID,
		&i.ShortName,
		&i.Description,
		&i.Url,
		&i.Color,
		&i.TextColor,
		&i.IsActive,
	)
	return i, err
}

const deleteServiceRoute = `-- name: DeleteServiceRoute :exec
DELETE FROM service_route
WHERE id = $1
`

func (q *Queries) DeleteServiceRoute(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteServiceRoute, id)
	return err
}

const deleteServiceRouteByServiceId = `-- name: DeleteServiceRouteByServiceId :exec
DELETE FROM service_route
WHERE service_id = $1
`

func (q *Queries) DeleteServiceRouteByServiceId(ctx context.Context, serviceID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteServiceRouteByServiceId, serviceID)
	return err
}

const getServiceRouteById = `-- name: GetServiceRouteById :one
SELECT id, created, modified, service_id, product_id, start_station_id, end_station_id, is_main_route, is_popular, ic_route_id, short_name, description, url, color, text_color, is_active FROM service_route
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetServiceRouteById(ctx context.Context, id pgtype.UUID) (ServiceRoute, error) {
	row := q.db.QueryRow(ctx, getServiceRouteById, id)
	var i ServiceRoute
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ServiceID,
		&i.ProductID,
		&i.StartStationID,
		&i.EndStationID,
		&i.IsMainRoute,
		&i.IsPopular,
		&i.IcRouteID,
		&i.ShortName,
		&i.Description,
		&i.Url,
		&i.Color,
		&i.TextColor,
		&i.IsActive,
	)
	return i, err
}

const getServiceRouteByServiceId = `-- name: GetServiceRouteByServiceId :one
SELECT id, created, modified, service_id, product_id, start_station_id, end_station_id, is_main_route, is_popular, ic_route_id, short_name, description, url, color, text_color, is_active FROM service_route
WHERE service_id = $1
LIMIT 1
`

func (q *Queries) GetServiceRouteByServiceId(ctx context.Context, serviceID pgtype.UUID) (ServiceRoute, error) {
	row := q.db.QueryRow(ctx, getServiceRouteByServiceId, serviceID)
	var i ServiceRoute
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ServiceID,
		&i.ProductID,
		&i.StartStationID,
		&i.EndStationID,
		&i.IsMainRoute,
		&i.IsPopular,
		&i.IcRouteID,
		&i.ShortName,
		&i.Description,
		&i.Url,
		&i.Color,
		&i.TextColor,
		&i.IsActive,
	)
	return i, err
}

const listServiceRoutesByServiceId = `-- name: ListServiceRoutesByServiceId :many
SELECT id, created, modified, service_id, product_id, start_station_id, end_station_id, is_main_route, is_popular, ic_route_id, short_name, description, url, color, text_color, is_active FROM service_route
WHERE service_id = $1
`

func (q *Queries) ListServiceRoutesByServiceId(ctx context.Context, serviceID pgtype.UUID) ([]ServiceRoute, error) {
	rows, err := q.db.Query(ctx, listServiceRoutesByServiceId, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ServiceRoute
	for rows.Next() {
		var i ServiceRoute
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.ServiceID,
			&i.ProductID,
			&i.StartStationID,
			&i.EndStationID,
			&i.IsMainRoute,
			&i.IsPopular,
			&i.IcRouteID,
			&i.ShortName,
			&i.Description,
			&i.Url,
			&i.Color,
			&i.TextColor,
			&i.IsActive,
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

const updateServiceRoute = `-- name: UpdateServiceRoute :one
UPDATE service_route
SET
    modified = now() AT TIME ZONE 'UTC',
    service_id = $2,
    product_id = $3,
    start_station_id = $4,
    end_station_id = $5,
    is_main_route = $6,
    is_popular = $7,
    ic_route_id = $8,
    short_name = $9,
    description = $10,
    url = $11,
    color = $12,
    text_color = $13
WHERE id = $1
RETURNING id, created, modified, service_id, product_id, start_station_id, end_station_id, is_main_route, is_popular, ic_route_id, short_name, description, url, color, text_color, is_active
`

type UpdateServiceRouteParams struct {
	ID             pgtype.UUID `db:"id" json:"id"`
	ServiceID      pgtype.UUID `db:"service_id" json:"serviceId"`
	ProductID      pgtype.UUID `db:"product_id" json:"productId"`
	StartStationID pgtype.UUID `db:"start_station_id" json:"startStationId"`
	EndStationID   pgtype.UUID `db:"end_station_id" json:"endStationId"`
	IsMainRoute    bool        `db:"is_main_route" json:"isMainRoute"`
	IsPopular      bool        `db:"is_popular" json:"isPopular"`
	IcRouteID      pgtype.Text `db:"ic_route_id" json:"icRouteId"`
	ShortName      pgtype.Text `db:"short_name" json:"shortName"`
	Description    pgtype.Text `db:"description" json:"description"`
	Url            pgtype.Text `db:"url" json:"url"`
	Color          pgtype.Text `db:"color" json:"color"`
	TextColor      pgtype.Text `db:"text_color" json:"textColor"`
}

func (q *Queries) UpdateServiceRoute(ctx context.Context, arg UpdateServiceRouteParams) (ServiceRoute, error) {
	row := q.db.QueryRow(ctx, updateServiceRoute,
		arg.ID,
		arg.ServiceID,
		arg.ProductID,
		arg.StartStationID,
		arg.EndStationID,
		arg.IsMainRoute,
		arg.IsPopular,
		arg.IcRouteID,
		arg.ShortName,
		arg.Description,
		arg.Url,
		arg.Color,
		arg.TextColor,
	)
	var i ServiceRoute
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ServiceID,
		&i.ProductID,
		&i.StartStationID,
		&i.EndStationID,
		&i.IsMainRoute,
		&i.IsPopular,
		&i.IcRouteID,
		&i.ShortName,
		&i.Description,
		&i.Url,
		&i.Color,
		&i.TextColor,
		&i.IsActive,
	)
	return i, err
}

const upsertServiceRoute = `-- name: UpsertServiceRoute :one
INSERT INTO service_route (
    id,
    created,
    modified,
    service_id,
    product_id,
    start_station_id,
    end_station_id,
    is_main_route,
    is_popular,
    ic_route_id,
    short_name,
    description,
    url,
    color,
    text_color
) VALUES (
             coalesce(nullif($13, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
         )
ON CONFLICT (id) DO UPDATE
    SET
        modified = now() AT TIME ZONE 'UTC',
        service_id = $1,
        product_id = $2,
        start_station_id = $3,
        end_station_id = $4,
        is_main_route = $5,
        is_popular = $6,
        ic_route_id = $7,
        short_name = $8,
        description = $9,
        url = $10,
        color = $11,
        text_color = $12
RETURNING id, created, modified, service_id, product_id, start_station_id, end_station_id, is_main_route, is_popular, ic_route_id, short_name, description, url, color, text_color, is_active
`

type UpsertServiceRouteParams struct {
	ServiceID      pgtype.UUID `db:"service_id" json:"serviceId"`
	ProductID      pgtype.UUID `db:"product_id" json:"productId"`
	StartStationID pgtype.UUID `db:"start_station_id" json:"startStationId"`
	EndStationID   pgtype.UUID `db:"end_station_id" json:"endStationId"`
	IsMainRoute    bool        `db:"is_main_route" json:"isMainRoute"`
	IsPopular      bool        `db:"is_popular" json:"isPopular"`
	IcRouteID      pgtype.Text `db:"ic_route_id" json:"icRouteId"`
	ShortName      pgtype.Text `db:"short_name" json:"shortName"`
	Description    pgtype.Text `db:"description" json:"description"`
	Url            pgtype.Text `db:"url" json:"url"`
	Color          pgtype.Text `db:"color" json:"color"`
	TextColor      pgtype.Text `db:"text_color" json:"textColor"`
	ID             interface{} `db:"id" json:"id"`
}

func (q *Queries) UpsertServiceRoute(ctx context.Context, arg UpsertServiceRouteParams) (ServiceRoute, error) {
	row := q.db.QueryRow(ctx, upsertServiceRoute,
		arg.ServiceID,
		arg.ProductID,
		arg.StartStationID,
		arg.EndStationID,
		arg.IsMainRoute,
		arg.IsPopular,
		arg.IcRouteID,
		arg.ShortName,
		arg.Description,
		arg.Url,
		arg.Color,
		arg.TextColor,
		arg.ID,
	)
	var i ServiceRoute
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ServiceID,
		&i.ProductID,
		&i.StartStationID,
		&i.EndStationID,
		&i.IsMainRoute,
		&i.IsPopular,
		&i.IcRouteID,
		&i.ShortName,
		&i.Description,
		&i.Url,
		&i.Color,
		&i.TextColor,
		&i.IsActive,
	)
	return i, err
}
