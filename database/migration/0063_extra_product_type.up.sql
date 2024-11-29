-- Extra Product Type Mapping Table
CREATE TABLE extra_product_type (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  extra_id uuid NOT NULL,
                                  product_type_id uuid NOT NULL,
                                  CONSTRAINT extra_id_fkey FOREIGN KEY (extra_id) REFERENCES extra
                                      ON DELETE CASCADE,
                                  CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type
                                      ON DELETE CASCADE
);

CREATE INDEX extra_product_type_extra_id_index ON extra_product_type USING HASH (extra_id);
CREATE INDEX extra_product_type_product_type_id_index ON extra_product_type USING HASH (product_type_id);
