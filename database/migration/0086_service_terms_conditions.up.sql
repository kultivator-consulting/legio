-- Service Terms & Conditions Table
CREATE TABLE service_terms_conditions (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    service_id uuid NOT NULL,
    ordering int NOT NULL,
    general_terms text,
    luggage_notes text,
    additional_notes text,
    CONSTRAINT service_id_fkey FOREIGN KEY (service_id) REFERENCES service
        ON DELETE CASCADE
);

CREATE INDEX service_terms_conditions_service_id_index ON service_terms_conditions USING HASH (service_id);
