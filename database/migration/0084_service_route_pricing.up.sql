-- Service Route Pricing Table
CREATE TABLE service_route_pricing (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    service_route_availability_id uuid NOT NULL,
    fare_type_id uuid NOT NULL,
    price_type varchar(30) NOT NULL DEFAULT 'One Way', -- ENUM 'One Way', 'Return'
    adult_price decimal NOT NULL,
    child_price decimal DEFAULT 0.00,
    infant_price decimal DEFAULT 0.00,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    excluded_dates timestamp with time zone[],
    sold_out_dates timestamp with time zone[],
    specific_terms text NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT service_route_availability_id_fkey FOREIGN KEY (service_route_availability_id) REFERENCES service_route_availability
        ON DELETE CASCADE,
    CONSTRAINT service_fare_type_id_fkey FOREIGN KEY (fare_type_id) REFERENCES fare_type
        ON DELETE CASCADE
);

CREATE INDEX service_route_pricing_service_route_availability_id_index ON service_route_pricing USING HASH (service_route_availability_id);
CREATE INDEX service_route_pricing_fare_type_id_index ON service_route_pricing USING HASH (fare_type_id);
