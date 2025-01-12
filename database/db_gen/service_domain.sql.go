// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: service_domain.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countServiceDomainsByDomainId = `-- name: CountServiceDomainsByDomainId :one
SELECT
    COUNT(id)
FROM service_domain
WHERE domain_id = $1
`

func (q *Queries) CountServiceDomainsByDomainId(ctx context.Context, domainID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countServiceDomainsByDomainId, domainID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countServiceDomainsByServiceId = `-- name: CountServiceDomainsByServiceId :one
SELECT
    COUNT(id)
FROM service_domain
WHERE service_id = $1
`

func (q *Queries) CountServiceDomainsByServiceId(ctx context.Context, serviceID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countServiceDomainsByServiceId, serviceID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createServiceDomain = `-- name: CreateServiceDomain :one
INSERT INTO service_domain (
    id,
    created,
    modified,
    service_id,
    domain_id
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             $1, $2
         )
RETURNING id, created, modified, service_id, domain_id
`

type CreateServiceDomainParams struct {
	ServiceID pgtype.UUID `db:"service_id" json:"serviceId"`
	DomainID  pgtype.UUID `db:"domain_id" json:"domainId"`
}

func (q *Queries) CreateServiceDomain(ctx context.Context, arg CreateServiceDomainParams) (ServiceDomain, error) {
	row := q.db.QueryRow(ctx, createServiceDomain, arg.ServiceID, arg.DomainID)
	var i ServiceDomain
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ServiceID,
		&i.DomainID,
	)
	return i, err
}

const deleteServiceDomain = `-- name: DeleteServiceDomain :exec
DELETE FROM service_domain
WHERE id = $1
`

func (q *Queries) DeleteServiceDomain(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteServiceDomain, id)
	return err
}

const deleteServiceDomainByServiceId = `-- name: DeleteServiceDomainByServiceId :exec
DELETE FROM service_domain
WHERE service_id = $1
`

func (q *Queries) DeleteServiceDomainByServiceId(ctx context.Context, serviceID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteServiceDomainByServiceId, serviceID)
	return err
}

const getServiceDomainByDomainId = `-- name: GetServiceDomainByDomainId :one
SELECT id, created, modified, service_id, domain_id FROM service_domain
WHERE domain_id = $1
LIMIT 1
`

func (q *Queries) GetServiceDomainByDomainId(ctx context.Context, domainID pgtype.UUID) (ServiceDomain, error) {
	row := q.db.QueryRow(ctx, getServiceDomainByDomainId, domainID)
	var i ServiceDomain
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ServiceID,
		&i.DomainID,
	)
	return i, err
}

const getServiceDomainById = `-- name: GetServiceDomainById :one
SELECT id, created, modified, service_id, domain_id FROM service_domain
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetServiceDomainById(ctx context.Context, id pgtype.UUID) (ServiceDomain, error) {
	row := q.db.QueryRow(ctx, getServiceDomainById, id)
	var i ServiceDomain
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ServiceID,
		&i.DomainID,
	)
	return i, err
}

const listServiceDomainsByDomainId = `-- name: ListServiceDomainsByDomainId :many
SELECT id, created, modified, service_id, domain_id FROM service_domain
WHERE domain_id = $1
`

func (q *Queries) ListServiceDomainsByDomainId(ctx context.Context, domainID pgtype.UUID) ([]ServiceDomain, error) {
	rows, err := q.db.Query(ctx, listServiceDomainsByDomainId, domainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ServiceDomain
	for rows.Next() {
		var i ServiceDomain
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.ServiceID,
			&i.DomainID,
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

const listServiceDomainsByServiceId = `-- name: ListServiceDomainsByServiceId :many
SELECT id, created, modified, service_id, domain_id FROM service_domain
WHERE service_id = $1
`

func (q *Queries) ListServiceDomainsByServiceId(ctx context.Context, serviceID pgtype.UUID) ([]ServiceDomain, error) {
	rows, err := q.db.Query(ctx, listServiceDomainsByServiceId, serviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ServiceDomain
	for rows.Next() {
		var i ServiceDomain
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.ServiceID,
			&i.DomainID,
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

const updateServiceDomain = `-- name: UpdateServiceDomain :one
UPDATE service_domain
SET
    modified = now() AT TIME ZONE 'UTC',
    service_id = $2,
    domain_id = $3
WHERE id = $1
RETURNING id, created, modified, service_id, domain_id
`

type UpdateServiceDomainParams struct {
	ID        pgtype.UUID `db:"id" json:"id"`
	ServiceID pgtype.UUID `db:"service_id" json:"serviceId"`
	DomainID  pgtype.UUID `db:"domain_id" json:"domainId"`
}

func (q *Queries) UpdateServiceDomain(ctx context.Context, arg UpdateServiceDomainParams) (ServiceDomain, error) {
	row := q.db.QueryRow(ctx, updateServiceDomain, arg.ID, arg.ServiceID, arg.DomainID)
	var i ServiceDomain
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.ServiceID,
		&i.DomainID,
	)
	return i, err
}
