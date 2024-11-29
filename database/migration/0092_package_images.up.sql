-- Package Images Mapping Table
CREATE TABLE package_images (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  package_id uuid NOT NULL,
                                  ordering integer NOT NULL,
                                  image varchar(400) NOT NULL,
                                  image_info varchar(600) NOT NULL,
                                  CONSTRAINT package_id_fkey FOREIGN KEY (package_id) REFERENCES package
                                      ON DELETE CASCADE
);

CREATE INDEX package_images_package_id_index ON package_images USING HASH (package_id);

