-- Package Availability Table
CREATE TABLE package_availability (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  package_id uuid NOT NULL,
                                  frequency varchar(40) NOT NULL,
                                  start_date timestamp with time zone NOT NULL,
                                  end_date timestamp with time zone NOT NULL,
                                  excluded_dates timestamp with time zone[] NOT NULL,
                                  sold_out_dates timestamp with time zone[] NOT NULL,
                                  CONSTRAINT package_id_fkey FOREIGN KEY (package_id) REFERENCES package
                                      ON DELETE CASCADE
);

CREATE INDEX package_availability_package_id_index ON package_availability USING HASH (package_id);

