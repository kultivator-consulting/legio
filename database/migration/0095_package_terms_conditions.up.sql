-- Package Terms & Conditions Table
CREATE TABLE package_terms_conditions (
                                        id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        package_id uuid NOT NULL,
                                        general_terms text,
                                        luggage_notes text,
                                        additional_notes text,
                                        CONSTRAINT package_id_fkey FOREIGN KEY (package_id) REFERENCES package
                                            ON DELETE CASCADE
);

CREATE INDEX package_terms_conditions_package_id_index ON package_terms_conditions USING HASH (package_id);
