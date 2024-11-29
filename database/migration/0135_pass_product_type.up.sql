-- Pass Product Type Mapping Table
CREATE TABLE pass_product_type (
      id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
      created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
      modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
      pass_id uuid NOT NULL,
      product_type_id uuid NOT NULL,
      category varchar(120) NOT NULL,
      CONSTRAINT pass_id_fkey FOREIGN KEY (pass_id) REFERENCES pass
          ON DELETE CASCADE,
      CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type
          ON DELETE CASCADE
);

CREATE INDEX pass_product_type_pass_id_index ON pass_product_type USING HASH (pass_id);
CREATE INDEX pass_product_type_product_type_id_index ON pass_product_type USING HASH (product_type_id);
