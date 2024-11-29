-- Pass Domain Mapping Table
CREATE TABLE pass_domain (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    pass_id uuid NOT NULL,
    domain_id uuid NOT NULL,
    CONSTRAINT pass_id_fkey FOREIGN KEY (pass_id) REFERENCES pass
        ON DELETE CASCADE,
    CONSTRAINT domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE CASCADE
);

CREATE INDEX pass_domain_pass_id_index ON pass_domain USING HASH (pass_id);
CREATE INDEX pass_domain_domain_id_index ON pass_domain USING HASH (domain_id);
