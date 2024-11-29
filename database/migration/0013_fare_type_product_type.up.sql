-- Fare Type Product Type Mapping Table
CREATE TABLE fare_type_product_type (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    fare_type_id uuid NOT NULL,
    product_type_id uuid NOT NULL,
    CONSTRAINT fare_type_id_fkey FOREIGN KEY (fare_type_id) REFERENCES fare_type
        ON DELETE CASCADE,
    CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type
        ON DELETE CASCADE
);

CREATE INDEX fare_type_product_type_fare_type_id_index ON fare_type_product_type USING HASH (fare_type_id);
CREATE INDEX fare_type_product_type_product_type_id_index ON fare_type_product_type USING HASH (product_type_id);
