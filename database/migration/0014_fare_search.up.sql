-- Fare Search Table
CREATE TABLE fare_search (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    product_type_id uuid NOT NULL,
    route_start_location_id uuid NOT NULL,
    route_end_location_id uuid NOT NULL,
    travel_date timestamp with time zone NOT NULL,
    is_return boolean NOT NULL DEFAULT FALSE,
    adult_count integer NOT NULL DEFAULT 0,
    child_count integer NOT NULL DEFAULT 0,
    infant_count integer NOT NULL DEFAULT 0,
    CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type
        ON DELETE CASCADE
);

CREATE INDEX fare_search_product_type_id_index ON fare_search USING HASH (product_type_id);
