-- Services Route Table
CREATE TABLE service_route (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    service_id uuid NOT NULL,
    product_id uuid DEFAULT NULL,
    start_station_id uuid NOT NULL,
    end_station_id uuid NOT NULL,
    is_main_route boolean NOT NULL DEFAULT false,
    is_popular boolean NOT NULL DEFAULT true,
    ic_route_id varchar(20) DEFAULT NULL,
    short_name varchar(300) DEFAULT NULL,
    description text,
    url varchar(400) DEFAULT NULL,
    color varchar(20) DEFAULT NULL,
    text_color varchar(20) DEFAULT NULL,
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT service_id_fkey FOREIGN KEY (service_id) REFERENCES service
        ON DELETE CASCADE,
    CONSTRAINT product_id_fkey FOREIGN KEY (product_id) REFERENCES product
        ON DELETE CASCADE,
    CONSTRAINT start_station_fkey FOREIGN KEY (start_station_id) REFERENCES station
        ON DELETE CASCADE,
    CONSTRAINT end_station_fkey FOREIGN KEY (end_station_id) REFERENCES station
        ON DELETE CASCADE
);

CREATE INDEX service_route_service_id_index ON service_route USING HASH (service_id);
CREATE INDEX service_route_product_id_index ON service_route USING HASH (product_id);
CREATE INDEX service_route_start_station_id_index ON service_route USING HASH (start_station_id);
CREATE INDEX service_route_end_station_id_index ON service_route USING HASH (end_station_id);
CREATE INDEX service_route_short_name_index ON service_route USING HASH (short_name);
CREATE INDEX service_route_ic_route_id_index ON service_route USING HASH (ic_route_id);
