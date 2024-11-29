-- Fare Search Result Availability Table
CREATE TABLE fare_search_result_availability (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    fare_search_result_id uuid NOT NULL,
    token varchar(255) DEFAULT NULL,
    fare_date timestamp with time zone NOT NULL,
	available boolean NOT NULL DEFAULT FALSE,
    availability_status text NOT NULL DEFAULT 'Unavailable',
    service_name text NOT NULL,
    departure_city text NOT NULL,
    departure_city_info text NOT NULL,
    arrival_city text NOT NULL,
    arrival_city_info text NOT NULL,
    departure_time timestamp with time zone NOT NULL,
    arrival_time timestamp with time zone NOT NULL,
    fare_type text NOT NULL,
    fare_type_description text NOT NULL,
    adult_price decimal NOT NULL,
    child_price decimal NOT NULL,
    infant_price decimal NOT NULL,
    description text NOT NULL,
    ticket_conditions text NOT NULL,
    CONSTRAINT fare_search_result_id_fkey FOREIGN KEY (fare_search_result_id) REFERENCES fare_search_result
        ON DELETE CASCADE
);

CREATE INDEX fare_search_result_availability_fare_search_result_id_index ON fare_search_result_availability USING HASH (fare_search_result_id);
