-- Service Domain Mapping Table
CREATE TABLE service_domain (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    service_id uuid NOT NULL,
    domain_id uuid NOT NULL,
    CONSTRAINT service_id_fkey FOREIGN KEY (service_id) REFERENCES service
        ON DELETE CASCADE,
    CONSTRAINT domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE CASCADE
);

CREATE INDEX service_domain_service_id_index ON service_domain USING HASH (service_id);
CREATE INDEX service_domain_domain_id_index ON service_domain USING HASH (domain_id);
