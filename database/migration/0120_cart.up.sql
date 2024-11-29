-- Cart Table
CREATE TABLE cart (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    cart_session_id varchar(120) NOT NULL,
    domain_id uuid NOT NULL,
    completed boolean NOT NULL DEFAULT false,
    CONSTRAINT cart_domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE RESTRICT
);

CREATE INDEX cart_session_id_index ON cart USING HASH (cart_session_id);
CREATE INDEX cart_domain_id_index ON cart USING HASH (domain_id);

