-- Package Itinerary Table
CREATE TABLE package_itinerary (
                                        id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        package_id uuid NOT NULL,
                                        day integer NOT NULL,
                                        title varchar(120) NOT NULL,
                                        event_icon varchar(120) NOT NULL,
                                        station_id uuid,
                                        latitude double precision DEFAULT NULL,
                                        longitude double precision DEFAULT NULL,
                                        description text NOT NULL,
                                        supplier varchar(120),
                                        supplier_code varchar(120),
                                        CONSTRAINT package_id_fkey FOREIGN KEY (package_id) REFERENCES package
                                            ON DELETE CASCADE,
                                        CONSTRAINT station_id_fkey FOREIGN KEY (station_id) REFERENCES station
                                            ON DELETE CASCADE
);

CREATE INDEX package_itinerary_package_id_index ON package_itinerary USING HASH (package_id);
CREATE INDEX package_itinerary_station_id_index ON package_itinerary USING HASH (station_id);
