-- Service Note Table
CREATE TABLE service_note (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    service_id uuid NOT NULL,
    note text NOT NULL,
    CONSTRAINT service_id_fkey FOREIGN KEY (service_id) REFERENCES service
        ON DELETE CASCADE
);

CREATE INDEX service_note_service_id_index ON service_note USING HASH (service_id);

