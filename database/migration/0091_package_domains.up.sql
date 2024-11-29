-- Package Domain Mapping Table
CREATE TABLE package_domains (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  package_id uuid NOT NULL,
                                  domain_id uuid NOT NULL,
                                  CONSTRAINT package_id_fkey FOREIGN KEY (package_id) REFERENCES package
                                      ON DELETE CASCADE,
                                  CONSTRAINT domain_id_fkey FOREIGN KEY (domain_id) REFERENCES domain
                                      ON DELETE CASCADE
);

CREATE INDEX package_domains_package_id_index ON package_domains USING HASH (package_id);
CREATE INDEX package_domains_domain_id_index ON package_domains USING HASH (domain_id);
