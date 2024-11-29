-- Product Table
CREATE TABLE product (
                    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
                    product_type_id uuid NOT NULL,
                    operator_id uuid DEFAULT NULL,
                    start_location_id uuid DEFAULT NULL,
                    end_location_id uuid DEFAULT NULL,
                    name varchar(300) NOT NULL,
                    start_place varchar(300) DEFAULT NULL,
                    end_place varchar(300) DEFAULT NULL,
                    operator_code varchar(300) DEFAULT NULL,
                    instructions text,
                    notes text,
                    is_active boolean DEFAULT true,
                    CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type,
                    CONSTRAINT operator_id_fkey FOREIGN KEY (operator_id) REFERENCES operator,
                    CONSTRAINT start_location_id_fkey FOREIGN KEY (start_location_id) REFERENCES location,
                    CONSTRAINT end_location_id_fkey FOREIGN KEY (end_location_id) REFERENCES location
);

CREATE INDEX product_name_index ON product USING HASH (name);
CREATE INDEX product_operator_code_index ON product USING HASH (operator_code);
CREATE INDEX product_product_type_id_index ON product USING HASH (product_type_id);
CREATE INDEX product_operator_id_index ON product USING HASH (operator_id);
CREATE INDEX product_start_location_id_index ON product USING HASH (start_location_id);
CREATE INDEX product_end_location_id_index ON product USING HASH (end_location_id);
