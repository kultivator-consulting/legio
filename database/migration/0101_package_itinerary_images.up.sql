-- Package Itinerary Images Mapping Table
CREATE TABLE package_itinerary_images (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  package_itinerary_id uuid NOT NULL,
                                  ordering integer NOT NULL,
                                  image varchar(400) NOT NULL,
                                  image_info varchar(600) NOT NULL,
                                  is_active boolean NOT NULL DEFAULT TRUE,
                                  CONSTRAINT package_id_fkey FOREIGN KEY (package_itinerary_id) REFERENCES package_itinerary
                                      ON DELETE CASCADE
);

CREATE INDEX package_itinerary_images_package_id_index ON package_itinerary_images USING HASH (package_itinerary_id);

