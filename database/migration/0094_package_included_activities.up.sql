-- Package Included Activities Table
CREATE TABLE package_included_activities (
                                        id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        package_id uuid NOT NULL,
                                        ordering int NOT NULL,
                                        title varchar(300) NOT NULL,
                                        description text,
                                        hero_image varchar(400) NOT NULL,
                                        hero_image_info varchar(600) NOT NULL,
                                        is_active boolean NOT NULL DEFAULT TRUE,
                                        CONSTRAINT package_id_fkey FOREIGN KEY (package_id) REFERENCES package
                                            ON DELETE CASCADE
);

CREATE INDEX package_included_activities_package_id_index ON package_included_activities USING HASH (package_id);
