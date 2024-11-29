-- Fare Search Result Table
CREATE TABLE fare_search_result (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    fare_search_id uuid NOT NULL,
    product_type_id uuid NOT NULL,
    route_start_location varchar(300) NOT NULL,
    route_end_location varchar(300) NOT NULL,
    fare_date timestamp with time zone NOT NULL,
    is_return boolean NOT NULL DEFAULT FALSE,
    CONSTRAINT fare_search_id_fkey FOREIGN KEY (fare_search_id) REFERENCES fare_search
        ON DELETE CASCADE,
    CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type
        ON DELETE CASCADE
);

CREATE INDEX fare_search_result_fare_search_id_index ON fare_search_result USING HASH (fare_search_id);
CREATE INDEX fare_search_result_product_type_id_index ON fare_search_result USING HASH (product_type_id);
