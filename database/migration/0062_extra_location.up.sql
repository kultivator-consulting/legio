-- Extra Location Mapping Table
CREATE TABLE extra_location (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  extra_id uuid NOT NULL,
                                  location_id uuid NOT NULL,
                                  CONSTRAINT extra_id_fkey FOREIGN KEY (extra_id) REFERENCES extra
                                      ON DELETE CASCADE,
                                  CONSTRAINT location_id_fkey FOREIGN KEY (location_id) REFERENCES location
                                      ON DELETE CASCADE
);

CREATE INDEX extra_location_extra_id_index ON extra_location USING HASH (extra_id);
CREATE INDEX extra_location_location_id_index ON extra_location USING HASH (location_id);
