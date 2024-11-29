-- Service Route Availability Table
CREATE TABLE service_route_availability (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    service_route_id uuid NOT NULL,
    departure_time timestamp with time zone NOT NULL,
    arrival_time timestamp with time zone NOT NULL,
    frequency varchar(40) NOT NULL,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    excluded_dates timestamp with time zone[] NOT NULL,
    sold_out_dates timestamp with time zone[] NOT NULL,
    ci_time_id varchar(200) DEFAULT NULL,
    notes text,
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT service_route_id_fkey FOREIGN KEY (service_route_id) REFERENCES service_route
        ON DELETE CASCADE
);

CREATE INDEX service_route_availability_service_route_id_index ON service_route_availability USING HASH (service_route_id);
CREATE INDEX service_route_availability_start_date_index ON service_route_availability USING HASH (start_date);
CREATE INDEX service_route_availability_end_date_index ON service_route_availability USING HASH (end_date);
CREATE INDEX service_route_availability_ci_time_id_index ON service_route_availability USING HASH (ci_time_id);
