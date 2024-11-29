-- Pass Availability Pricing Table
CREATE TABLE pass_availability_pricing (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    pass_availability_id uuid NOT NULL,
    fare_type_id uuid NOT NULL,
    adult_price decimal NOT NULL,
    child_price decimal DEFAULT 0.00,
    infant_price decimal DEFAULT 0.00,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    excluded_dates timestamp with time zone[],
    CONSTRAINT pass_availability_pass_availability_id_fkey FOREIGN KEY (pass_availability_id) REFERENCES pass_availability
        ON DELETE CASCADE,
    CONSTRAINT pass_availability_pricing_fare_type_id_fkey FOREIGN KEY (fare_type_id) REFERENCES fare_type
        ON DELETE CASCADE
);

CREATE INDEX pass_availability_pass_availability_id_index ON pass_availability_pricing USING HASH (pass_availability_id);
CREATE INDEX pass_availability_pricing_fare_type_id_index ON pass_availability_pricing USING HASH (fare_type_id);
