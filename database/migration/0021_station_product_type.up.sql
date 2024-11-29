-- Station Product Type Mapping Table
CREATE TABLE station_product_type (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    station_id uuid NOT NULL,
    product_type_id uuid NOT NULL,
    CONSTRAINT station_id_fkey FOREIGN KEY (station_id) REFERENCES station
        ON DELETE CASCADE,
    CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type
        ON DELETE NO ACTION
);

CREATE INDEX station_product_type_station_id_index ON station_product_type USING HASH (station_id);
CREATE INDEX station_product_type_product_type_id_index ON station_product_type USING HASH (product_type_id);
