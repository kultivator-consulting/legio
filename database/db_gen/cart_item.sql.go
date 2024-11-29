// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: cart_item.sql

package db_gen

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countCartItems = `-- name: CountCartItems :one
SELECT
    COUNT(id)
FROM cart_item
WHERE deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountCartItems(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countCartItems)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countCartItemsByCartId = `-- name: CountCartItemsByCartId :one
SELECT
    COUNT(id)
FROM cart_item
WHERE cart_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
`

func (q *Queries) CountCartItemsByCartId(ctx context.Context, cartID pgtype.UUID) (int64, error) {
	row := q.db.QueryRow(ctx, countCartItemsByCartId, cartID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createCartItem = `-- name: CreateCartItem :one
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
RETURNING id, created, modified, deleted, cart_id, associated_cart_item_id, item_id, item_type, line_code, description, quantity, price, discount, data, rule_handler, is_discount
`

type CreateCartItemParams struct {
	CartID               pgtype.UUID    `db:"cart_id" json:"cartId"`
	AssociatedCartItemID pgtype.UUID    `db:"associated_cart_item_id" json:"associatedCartItemId"`
	ItemID               pgtype.UUID    `db:"item_id" json:"itemId"`
	ItemType             string         `db:"item_type" json:"itemType"`
	LineCode             string         `db:"line_code" json:"lineCode"`
	Description          string         `db:"description" json:"description"`
	Quantity             pgtype.Numeric `db:"quantity" json:"quantity"`
	Price                pgtype.Numeric `db:"price" json:"price"`
	Discount             pgtype.Numeric `db:"discount" json:"discount"`
	Data                 string         `db:"data" json:"data"`
	RuleHandler          string         `db:"rule_handler" json:"ruleHandler"`
	IsDiscount           bool           `db:"is_discount" json:"isDiscount"`
}

func (q *Queries) CreateCartItem(ctx context.Context, arg CreateCartItemParams) (CartItem, error) {
	row := q.db.QueryRow(ctx, createCartItem,
		arg.CartID,
		arg.AssociatedCartItemID,
		arg.ItemID,
		arg.ItemType,
		arg.LineCode,
		arg.Description,
		arg.Quantity,
		arg.Price,
		arg.Discount,
		arg.Data,
		arg.RuleHandler,
		arg.IsDiscount,
	)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.CartID,
		&i.AssociatedCartItemID,
		&i.ItemID,
		&i.ItemType,
		&i.LineCode,
		&i.Description,
		&i.Quantity,
		&i.Price,
		&i.Discount,
		&i.Data,
		&i.RuleHandler,
		&i.IsDiscount,
	)
	return i, err
}

const deleteCartItemById = `-- name: DeleteCartItemById :exec
UPDATE cart_item
SET
    deleted = now() AT TIME ZONE 'UTC'
WHERE id = $1
`

func (q *Queries) DeleteCartItemById(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteCartItemById, id)
	return err
}

const getCartItemByCartId = `-- name: GetCartItemByCartId :one
SELECT id, created, modified, deleted, cart_id, associated_cart_item_id, item_id, item_type, line_code, description, quantity, price, discount, data, rule_handler, is_discount FROM cart_item
WHERE cart_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetCartItemByCartId(ctx context.Context, cartID pgtype.UUID) (CartItem, error) {
	row := q.db.QueryRow(ctx, getCartItemByCartId, cartID)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.CartID,
		&i.AssociatedCartItemID,
		&i.ItemID,
		&i.ItemType,
		&i.LineCode,
		&i.Description,
		&i.Quantity,
		&i.Price,
		&i.Discount,
		&i.Data,
		&i.RuleHandler,
		&i.IsDiscount,
	)
	return i, err
}

const getCartItemById = `-- name: GetCartItemById :one
SELECT id, created, modified, deleted, cart_id, associated_cart_item_id, item_id, item_type, line_code, description, quantity, price, discount, data, rule_handler, is_discount FROM cart_item
WHERE id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
LIMIT 1
`

func (q *Queries) GetCartItemById(ctx context.Context, id pgtype.UUID) (CartItem, error) {
	row := q.db.QueryRow(ctx, getCartItemById, id)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.CartID,
		&i.AssociatedCartItemID,
		&i.ItemID,
		&i.ItemType,
		&i.LineCode,
		&i.Description,
		&i.Quantity,
		&i.Price,
		&i.Discount,
		&i.Data,
		&i.RuleHandler,
		&i.IsDiscount,
	)
	return i, err
}

const listAllCartItems = `-- name: ListAllCartItems :many
SELECT id, created, modified, deleted, cart_id, associated_cart_item_id, item_id, item_type, line_code, description, quantity, price, discount, data, rule_handler, is_discount FROM cart_item
WHERE deleted > now() AT TIME ZONE 'UTC'
ORDER BY created
`

func (q *Queries) ListAllCartItems(ctx context.Context) ([]CartItem, error) {
	rows, err := q.db.Query(ctx, listAllCartItems)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CartItem
	for rows.Next() {
		var i CartItem
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.CartID,
			&i.AssociatedCartItemID,
			&i.ItemID,
			&i.ItemType,
			&i.LineCode,
			&i.Description,
			&i.Quantity,
			&i.Price,
			&i.Discount,
			&i.Data,
			&i.RuleHandler,
			&i.IsDiscount,
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

const listCartItemsByCartId = `-- name: ListCartItemsByCartId :many
SELECT id, created, modified, deleted, cart_id, associated_cart_item_id, item_id, item_type, line_code, description, quantity, price, discount, data, rule_handler, is_discount FROM cart_item
WHERE cart_id = $1
  AND deleted > now() AT TIME ZONE 'UTC'
ORDER BY created
`

func (q *Queries) ListCartItemsByCartId(ctx context.Context, cartID pgtype.UUID) ([]CartItem, error) {
	rows, err := q.db.Query(ctx, listCartItemsByCartId, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CartItem
	for rows.Next() {
		var i CartItem
		if err := rows.Scan(
			&i.ID,
			&i.Created,
			&i.Modified,
			&i.Deleted,
			&i.CartID,
			&i.AssociatedCartItemID,
			&i.ItemID,
			&i.ItemType,
			&i.LineCode,
			&i.Description,
			&i.Quantity,
			&i.Price,
			&i.Discount,
			&i.Data,
			&i.RuleHandler,
			&i.IsDiscount,
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

const updateCartItemById = `-- name: UpdateCartItemById :one
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
RETURNING id, created, modified, deleted, cart_id, associated_cart_item_id, item_id, item_type, line_code, description, quantity, price, discount, data, rule_handler, is_discount
`

type UpdateCartItemByIdParams struct {
	ID                   pgtype.UUID    `db:"id" json:"id"`
	CartID               pgtype.UUID    `db:"cart_id" json:"cartId"`
	AssociatedCartItemID pgtype.UUID    `db:"associated_cart_item_id" json:"associatedCartItemId"`
	ItemID               pgtype.UUID    `db:"item_id" json:"itemId"`
	ItemType             string         `db:"item_type" json:"itemType"`
	LineCode             string         `db:"line_code" json:"lineCode"`
	Description          string         `db:"description" json:"description"`
	Quantity             pgtype.Numeric `db:"quantity" json:"quantity"`
	Price                pgtype.Numeric `db:"price" json:"price"`
	Discount             pgtype.Numeric `db:"discount" json:"discount"`
	Data                 string         `db:"data" json:"data"`
	RuleHandler          string         `db:"rule_handler" json:"ruleHandler"`
	IsDiscount           bool           `db:"is_discount" json:"isDiscount"`
}

func (q *Queries) UpdateCartItemById(ctx context.Context, arg UpdateCartItemByIdParams) (CartItem, error) {
	row := q.db.QueryRow(ctx, updateCartItemById,
		arg.ID,
		arg.CartID,
		arg.AssociatedCartItemID,
		arg.ItemID,
		arg.ItemType,
		arg.LineCode,
		arg.Description,
		arg.Quantity,
		arg.Price,
		arg.Discount,
		arg.Data,
		arg.RuleHandler,
		arg.IsDiscount,
	)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.CartID,
		&i.AssociatedCartItemID,
		&i.ItemID,
		&i.ItemType,
		&i.LineCode,
		&i.Description,
		&i.Quantity,
		&i.Price,
		&i.Discount,
		&i.Data,
		&i.RuleHandler,
		&i.IsDiscount,
	)
	return i, err
}

const upsertCartItem = `-- name: UpsertCartItem :one
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
             coalesce(nullif($13, uuid_nil()), uuid_generate_v4()),
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
RETURNING id, created, modified, deleted, cart_id, associated_cart_item_id, item_id, item_type, line_code, description, quantity, price, discount, data, rule_handler, is_discount
`

type UpsertCartItemParams struct {
	CartID               pgtype.UUID    `db:"cart_id" json:"cartId"`
	AssociatedCartItemID pgtype.UUID    `db:"associated_cart_item_id" json:"associatedCartItemId"`
	ItemID               pgtype.UUID    `db:"item_id" json:"itemId"`
	ItemType             string         `db:"item_type" json:"itemType"`
	LineCode             string         `db:"line_code" json:"lineCode"`
	Description          string         `db:"description" json:"description"`
	Quantity             pgtype.Numeric `db:"quantity" json:"quantity"`
	Price                pgtype.Numeric `db:"price" json:"price"`
	Discount             pgtype.Numeric `db:"discount" json:"discount"`
	Data                 string         `db:"data" json:"data"`
	RuleHandler          string         `db:"rule_handler" json:"ruleHandler"`
	IsDiscount           bool           `db:"is_discount" json:"isDiscount"`
	ID                   interface{}    `db:"id" json:"id"`
}

func (q *Queries) UpsertCartItem(ctx context.Context, arg UpsertCartItemParams) (CartItem, error) {
	row := q.db.QueryRow(ctx, upsertCartItem,
		arg.CartID,
		arg.AssociatedCartItemID,
		arg.ItemID,
		arg.ItemType,
		arg.LineCode,
		arg.Description,
		arg.Quantity,
		arg.Price,
		arg.Discount,
		arg.Data,
		arg.RuleHandler,
		arg.IsDiscount,
		arg.ID,
	)
	var i CartItem
	err := row.Scan(
		&i.ID,
		&i.Created,
		&i.Modified,
		&i.Deleted,
		&i.CartID,
		&i.AssociatedCartItemID,
		&i.ItemID,
		&i.ItemType,
		&i.LineCode,
		&i.Description,
		&i.Quantity,
		&i.Price,
		&i.Discount,
		&i.Data,
		&i.RuleHandler,
		&i.IsDiscount,
	)
	return i, err
}