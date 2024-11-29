// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: package_images.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countPackageImagesByPackageId = `-- name: CountPackageImagesByPackageId :one
SELECT
    COUNT(id)
FROM package_images
WHERE package_id = $1
`

func (q *Queries) CountPackageImagesByPackageId(ctx context.Context, packageID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countPackageImagesByPackageId, packageID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createPackageImage = `-- name: CreatePackageImage :one
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
RETURNING id, created, modified, package_id, ordering, image, image_info
`

type CreatePackageImageParams struct {
	PackageID pgtype.UUID `db:"package_id" json:"packageId"`
	Ordering  int32       `db:"ordering" json:"ordering"`
	Image     string      `db:"image" json:"image"`
	ImageInfo string      `db:"image_info" json:"imageInfo"`
}

func (q *Queries) CreatePackageImage(ctx context.Context, arg CreatePackageImageParams) (PackageImages, error) {
	row := q.db.QueryRow(ctx, createPackageImage,
		arg.PackageID,
		arg.Ordering,
		arg.Image,
		arg.ImageInfo,
	)
	var i PackageImages
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PackageID,
		&i.Ordering,
		&i.Image,
		&i.ImageInfo,
	)
	return i, err
}

const deletePackageImage = `-- name: DeletePackageImage :exec
DELETE FROM package_images
WHERE id = $1
`

func (q *Queries) DeletePackageImage(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePackageImage, id)
	return err
}

const deletePackageImagesByPackageId = `-- name: DeletePackageImagesByPackageId :exec
DELETE FROM package_images
WHERE package_id = $1
`

func (q *Queries) DeletePackageImagesByPackageId(ctx context.Context, packageID pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deletePackageImagesByPackageId, packageID)
	return err
}

const getPackageImageById = `-- name: GetPackageImageById :one
SELECT id, created, modified, package_id, ordering, image, image_info FROM package_images
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetPackageImageById(ctx context.Context, id pgtype.UUID) (PackageImages, error) {
	row := q.db.QueryRow(ctx, getPackageImageById, id)
	var i PackageImages
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PackageID,
		&i.Ordering,
		&i.Image,
		&i.ImageInfo,
	)
	return i, err
}

const getPackageImagesByPackageId = `-- name: GetPackageImagesByPackageId :many
SELECT id, created, modified, package_id, ordering, image, image_info FROM package_images
WHERE package_id = $1
`

func (q *Queries) GetPackageImagesByPackageId(ctx context.Context, packageID pgtype.UUID) ([]PackageImages, error) {
	rows, err := q.db.Query(ctx, getPackageImagesByPackageId, packageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PackageImages
	for rows.Next() {
		var i PackageImages
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.PackageID,
			&i.Ordering,
			&i.Image,
			&i.ImageInfo,
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

const listPackageImagesByPackageIdAsc = `-- name: ListPackageImagesByPackageIdAsc :many
SELECT id, created, modified, package_id, ordering, image, image_info FROM package_images
WHERE package_id = $1
ORDER BY $2::text
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPackageImagesByPackageIdAscParams struct {
	PackageID         pgtype.UUID `db:"package_id" json:"packageId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPackageImagesByPackageIdAsc(ctx context.Context, arg ListPackageImagesByPackageIdAscParams) ([]PackageImages, error) {
	rows, err := q.db.Query(ctx, listPackageImagesByPackageIdAsc,
		arg.PackageID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PackageImages
	for rows.Next() {
		var i PackageImages
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.PackageID,
			&i.Ordering,
			&i.Image,
			&i.ImageInfo,
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

const listPackageImagesByPackageIdDesc = `-- name: ListPackageImagesByPackageIdDesc :many
SELECT id, created, modified, package_id, ordering, image, image_info FROM package_images
WHERE package_id = $1
ORDER BY $2::text DESC
OFFSET ($3::int - 1) * $4::int
    FETCH NEXT $4 ROWS ONLY
`

type ListPackageImagesByPackageIdDescParams struct {
	PackageID         pgtype.UUID `db:"package_id" json:"packageId"`
	SortBy            string      `db:"sort_by" json:"sortBy"`
	RequestedPage     int32       `db:"requested_page" json:"requestedPage"`
	RequestedPageSize int32       `db:"requested_page_size" json:"requestedPageSize"`
}

func (q *Queries) ListPackageImagesByPackageIdDesc(ctx context.Context, arg ListPackageImagesByPackageIdDescParams) ([]PackageImages, error) {
	rows, err := q.db.Query(ctx, listPackageImagesByPackageIdDesc,
		arg.PackageID,
		arg.SortBy,
		arg.RequestedPage,
		arg.RequestedPageSize,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PackageImages
	for rows.Next() {
		var i PackageImages
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.PackageID,
			&i.Ordering,
			&i.Image,
			&i.ImageInfo,
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

const updatePackageImage = `-- name: UpdatePackageImage :one
UPDATE package_images
SET
    modified = now() AT TIME ZONE 'UTC',
    package_id = $2,
    ordering = $3,
    image = $4,
    image_info = $5
WHERE id = $1
RETURNING id, created, modified, package_id, ordering, image, image_info
`

type UpdatePackageImageParams struct {
	ID        pgtype.UUID `db:"id" json:"id"`
	PackageID pgtype.UUID `db:"package_id" json:"packageId"`
	Ordering  int32       `db:"ordering" json:"ordering"`
	Image     string      `db:"image" json:"image"`
	ImageInfo string      `db:"image_info" json:"imageInfo"`
}

func (q *Queries) UpdatePackageImage(ctx context.Context, arg UpdatePackageImageParams) (PackageImages, error) {
	row := q.db.QueryRow(ctx, updatePackageImage,
		arg.ID,
		arg.PackageID,
		arg.Ordering,
		arg.Image,
		arg.ImageInfo,
	)
	var i PackageImages
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.PackageID,
		&i.Ordering,
		&i.Image,
		&i.ImageInfo,
	)
	return i, err
}
