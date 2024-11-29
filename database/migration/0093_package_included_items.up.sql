-- Package Included Items Table
CREATE TABLE package_included_items (
                                        id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        package_id uuid NOT NULL,
                                        ordering int NOT NULL,
                                        type varchar(120) NOT NULL, -- Train Travel, Accommodation, Tours, Bus Travel, Meals, etc.
                                        description text,
                                        item_icon varchar(120) NOT NULL,
                                        is_active boolean NOT NULL DEFAULT TRUE,
                                        CONSTRAINT package_id_fkey FOREIGN KEY (package_id) REFERENCES package
                                            ON DELETE CASCADE
);

CREATE INDEX package_included_items_package_id_index ON package_included_items USING HASH (package_id);
