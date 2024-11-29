-- Operator Product Type Mapping Table
CREATE TABLE operator_product_type (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    operator_id uuid NOT NULL,
    product_type_id uuid NOT NULL,
    CONSTRAINT operator_id_fkey FOREIGN KEY (operator_id) REFERENCES operator
        ON DELETE CASCADE,
    CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type
        ON DELETE NO ACTION
);

CREATE INDEX operator_product_type_operator_id_index ON operator_product_type USING HASH (operator_id);
CREATE INDEX operator_product_type_product_type_id_index ON operator_product_type USING HASH (product_type_id);
