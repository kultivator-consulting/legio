-- name: GetCartItemById :one
SELECT * FROM cart_item
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: GetCartItemByCartId :one
SELECT * FROM cart_item
WHERE cart_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1;

-- name: CountCartItems :one
SELECT
    COUNT(id)
FROM cart_item
WHERE deleted > now() AT TIME ZONE 'UTC';

-- name: ListAllCartItems :many
SELECT * FROM cart_item
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY created;

-- name: CountCartItemsByCartId :one
SELECT
    COUNT(id)
FROM cart_item
WHERE cart_id = $1
  AND deleted > now() AT TIME ZONE 'UTC';

-- name: ListCartItemsByCartId :many
SELECT * FROM cart_item
WHERE cart_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY created;

-- name: CreateCartItem :one
INSERT INTO cart_item (
    id,
    created,
    modified,
    deleted,
    cart_id,
    associated_cart_item_id,
    item_id,
    item_type,
    line_code,
    description,
    quantity,
    price,
    discount,
    data,
    rule_handler,
    is_discount
) VALUES (
             uuid_generate_v4(),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
         )
RETURNING *;

-- name: UpsertCartItem :one
INSERT INTO cart_item (
    id,
    created,
    modified,
    deleted,
    cart_id,
    associated_cart_item_id,
    item_id,
    item_type,
    line_code,
    description,
    quantity,
    price,
    discount,
    data,
    rule_handler,
    is_discount
) VALUES (
             coalesce(nullif(@id, uuid_nil()), uuid_generate_v4()),
             now() AT TIME ZONE 'UTC',
             now() AT TIME ZONE 'UTC',
             'infinity'::timestamp AT TIME ZONE 'UTC',
             $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
         )
ON CONFLICT (id)
DO UPDATE SET
    modified = now() AT TIME ZONE 'UTC',
    cart_id = $1,
    associated_cart_item_id = $2,
    item_id = $3,
    item_type = $4,
    line_code = $5,
    description = $6,
    quantity = $7,
    price = $8,
    discount = $9,
    data = $10,
    rule_handler = $11,
    is_discount = $12
RETURNING *;

-- name: DeleteCartItemById :exec
UPDATE cart_item
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1;

-- name: UpdateCartItemById :one
UPDATE cart_item
SET
    modified = now() AT TIME ZONE 'UTC',
    cart_id = $2,
    associated_cart_item_id = $3,
    item_id = $4,
    item_type = $5,
    line_code = $6,
    description = $7,
    quantity = $8,
    price = $9,
    discount = $10,
    data = $11,
    rule_handler = $12,
    is_discount = $13
WHERE id = $1
RETURNING *;
