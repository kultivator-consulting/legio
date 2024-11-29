-- name: GetBlogById :one
SELECT * FROM blog
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetBlogByDomainIdAndId :one
SELECT * FROM blog
WHERE domain_id = $1
  AND id = $2
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetBlogByPageId :one
SELECT * FROM blog
WHERE page_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountBlogs :one
SELECT
    COUNT(id)
FROM blog
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListAllBlogs :many
SELECT * FROM blog
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY title;

-- name: CreateBlog :one
INSERT INTO blog (
    id,
    created,
    modified,
    deleted,
    domain_id,
    account_id,
    page_id,
    title,
    description,
    image,
    image_info,
    keywords,
    is_active
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9
         )
RETURNING *;

-- name: DeleteBlogById :exec
UPDATE blog
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateBlogById :one
UPDATE blog
SET
    modified = now() AT TIME ZONE 'UTC',
    domain_id = $2,
    account_id = $3,
    page_id = $4,
    title = $5,
    description = $6,
    image = $7,
    image_info = $8,
    keywords = $9,
    is_active = $10
WHERE id = $1
RETURNING *;
