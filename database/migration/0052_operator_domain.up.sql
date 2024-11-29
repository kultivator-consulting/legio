-- Operator Domain Mapping Table
CREATE TABLE operator_domain (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    operator_id uuid NOT NULL,
    domain_id uuid NOT NULL,
    CONSTRAINT operator_id_fkey FOREIGN KEY (operator_id) REFERENCES operator
        ON DELETE CASCADE,
    CONSTRAINT domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE NO ACTION
);

CREATE INDEX operator_domain_operator_id_index ON operator_domain USING HASH (operator_id);
CREATE INDEX operator_domain_domain_id_index ON operator_domain USING HASH (domain_id);
