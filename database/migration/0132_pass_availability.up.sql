-- Pass Availability Table
CREATE TABLE pass_availability (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    pass_id uuid NOT NULL,
    start_date timestamp with time zone NOT NULL,
    end_date timestamp with time zone NOT NULL,
    excluded_dates timestamp with time zone[] NOT NULL,
    notes text,
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT pass_availability_pass_id_fkey FOREIGN KEY (pass_id) REFERENCES pass
        ON DELETE CASCADE
);

CREATE INDEX pass_availability_pass_service_id_index ON pass_availability USING HASH (pass_id);
CREATE INDEX pass_availability_start_date_index ON pass_availability USING HASH (start_date);
CREATE INDEX pass_availability_end_date_index ON pass_availability USING HASH (end_date);
