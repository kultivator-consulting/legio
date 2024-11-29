-- Station Images Mapping Table
CREATE TABLE station_image (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  station_id uuid NOT NULL,
                                  ordering integer NOT NULL,
                                  image varchar(400) NOT NULL,
                                  image_info varchar(600) NOT NULL,
                                  CONSTRAINT station_id_fkey FOREIGN KEY (station_id) REFERENCES station
                                      ON DELETE CASCADE
);

CREATE INDEX station_image_station_id_index ON station_image USING HASH (station_id);

