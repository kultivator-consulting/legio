-- Station Table
CREATE TABLE station (
                         id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                         created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                         modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                         deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
                         location_id uuid DEFAULT NULL,
                         name varchar(300) NOT NULL,
                         address text NOT NULL,
                         code varchar(300) NOT NULL,
                         description text NOT NULL,
                         latitude varchar(300) NOT NULL,
                         longitude varchar(300) NOT NULL,
                         stop_id varchar(100) NOT NULL,
                         zone_id varchar(100) NOT NULL,
                         stop_url varchar(4000) NOT NULL,
                         is_active boolean NOT NULL DEFAULT true,
                         CONSTRAINT station_location_id_fkey FOREIGN KEY (location_id) REFERENCES location
                             ON DELETE CASCADE
);

CREATE INDEX station_location_id_index ON station USING HASH (location_id);
CREATE INDEX station_name_index ON station USING HASH (name);
CREATE INDEX station_code_index ON station USING HASH (code);
CREATE INDEX station_latitude_index ON station USING HASH (latitude);
CREATE INDEX station_is_active_index ON station USING HASH (is_active);
