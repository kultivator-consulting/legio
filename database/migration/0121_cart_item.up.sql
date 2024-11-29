-- Cart Item Table
CREATE TABLE cart_item (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    cart_id uuid NOT NULL,
    associated_cart_item_id uuid DEFAULT NULL, -- The cart item that this line item is associated with a discount line item, etc.
    item_id uuid NOT NULL,  -- This will be the id of the item, could be passes, bookings, etc.
    item_type varchar(20) NOT NULL,  -- This will be the type of the item, could be pass, booking, etc.
    line_code varchar(20) NOT NULL,
    description text NOT NULL,
    quantity decimal NOT NULL,
    price decimal NOT NULL,
    discount decimal NOT NULL DEFAULT 0,
    data text NOT NULL,
    rule_handler varchar(20) NOT NULL,
    is_discount boolean NOT NULL DEFAULT false,
    CONSTRAINT cart_id_fkey FOREIGN KEY (cart_id) REFERENCES cart
        ON DELETE CASCADE
);

CREATE INDEX cart_item_item_id_index ON cart_item USING HASH (item_id);
CREATE INDEX cart_item_associated_cart_item_id_index ON cart_item USING HASH (associated_cart_item_id);
CREATE INDEX cart_item_line_code_index ON cart_item USING HASH (line_code);

