-- Service Stop Ordering Table
CREATE TABLE service_stop_ordering (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    service_id uuid NOT NULL,
    location_id uuid NOT NULL,
    ordering integer NOT NULL,
    CONSTRAINT service_id_fkey FOREIGN KEY (service_id) REFERENCES service
        ON DELETE CASCADE,
    CONSTRAINT location_id_fkey FOREIGN KEY (location_id) REFERENCES location
        ON DELETE CASCADE
);

CREATE INDEX service_stop_ordering_service_id_index ON service_stop_ordering USING HASH (service_id);
CREATE INDEX service_stop_ordering_location_id_index ON service_stop_ordering USING HASH (location_id);
