-- Product Type Domain Mapping Table
CREATE TABLE product_type_domain (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    product_type_id uuid NOT NULL,
    domain_id uuid NOT NULL,
    CONSTRAINT product_type_id_fkey FOREIGN KEY (product_type_id) REFERENCES product_type
        ON DELETE CASCADE,
    CONSTRAINT domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
        ON DELETE NO ACTION
);

CREATE INDEX product_type_domain_product_type_id_index ON product_type_domain USING HASH (product_type_id);
CREATE INDEX product_type_domain_domain_id_index ON product_type_domain USING HASH (domain_id);
