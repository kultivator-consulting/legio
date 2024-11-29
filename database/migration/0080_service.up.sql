-- Services Table
CREATE TABLE service (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    operator_id uuid NOT NULL,
    product_type_id uuid NOT NULL,
    name varchar(250) NOT NULL,
    sub_service_name varchar(255) DEFAULT NULL,
    start_location_id uuid NOT NULL,
    end_location_id uuid NOT NULL,
    start_station_id uuid DEFAULT NULL,
    end_station_id uuid DEFAULT NULL,
    logo varchar(250) DEFAULT NULL,
    banner_image varchar(250) DEFAULT NULL,
    excluded_dates timestamp with time zone[] NOT NULL DEFAULT '{}',
    sold_out_dates timestamp with time zone[] NOT NULL DEFAULT '{}',
    is_popular boolean NOT NULL DEFAULT TRUE,
    with_pass boolean NOT NULL DEFAULT FALSE,
    is_tour boolean NOT NULL DEFAULT FALSE,
    tour_name varchar(255) DEFAULT NULL,
    description text,
    day_excursion boolean NOT NULL DEFAULT FALSE,
    multi_service boolean NOT NULL DEFAULT FALSE,
    excursion_text text,
    service_code varchar(200) DEFAULT NULL,
    meta_title text,
    meta_content text,
    is_active boolean NOT NULL DEFAULT TRUE,
    CONSTRAINT operator_id_fkey FOREIGN KEY (operator_id) REFERENCES operator
        ON DELETE CASCADE,
    CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type
        ON DELETE NO ACTION,
    CONSTRAINT start_location_id_fkey FOREIGN KEY (start_location_id) REFERENCES location
        ON DELETE CASCADE,
    CONSTRAINT end_location_id_fkey FOREIGN KEY (end_location_id) REFERENCES location
        ON DELETE CASCADE,
    CONSTRAINT start_station_id_fkey FOREIGN KEY (start_station_id) REFERENCES station
        ON DELETE CASCADE,
    CONSTRAINT end_station_id_fkey FOREIGN KEY (end_station_id) REFERENCES station
        ON DELETE CASCADE
);

CREATE INDEX service_operator_id_index ON service USING HASH (operator_id);
CREATE INDEX service_product_type_id_index ON service USING HASH (product_type_id);
CREATE INDEX service_name_index ON service USING HASH (name);
CREATE INDEX service_is_active_index ON service USING HASH (is_active);
CREATE INDEX service_start_location_id_index ON service USING HASH (start_location_id);
CREATE INDEX service_end_location_id_index ON service USING HASH (end_location_id);
CREATE INDEX service_start_station_id_index ON service USING HASH (start_station_id);
CREATE INDEX service_end_station_id_index ON service USING HASH (end_station_id);

