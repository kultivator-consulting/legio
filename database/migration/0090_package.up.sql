-- Package Table
CREATE TABLE package (
                                        id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                        created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                        deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
                                        title varchar(120) NOT NULL,
                                        slug varchar(120) NOT NULL,
                                        code varchar(120) NOT NULL,
                                        introduction text NOT NULL, -- This is the short description of the package
                                        hero_image varchar(400) NOT NULL,
                                        hero_image_info varchar(600) NOT NULL,
                                        duration varchar(120) NOT NULL,
                                        required_deposit decimal NOT NULL,
                                        departs varchar(400) NOT NULL,
                                        destinations varchar(400)[] NOT NULL, -- This is the list of 'Destinations' or 'Destination' visible to the user
                                        route varchar(400)[] NOT NULL,
                                        description text NOT NULL, -- This is the 'What awaits you...' section
                                        included_description text NOT NULL, -- This is the 'What included in this package' section
                                        activities_description text NOT NULL, -- This is the 'Included activities' section
                                        journey_map varchar(400) NOT NULL,
                                        island_filter varchar(200) NOT NULL,
                                        keywords text NOT NULL,
                                        is_group_journey boolean NOT NULL DEFAULT FALSE,
                                        is_active boolean NOT NULL DEFAULT TRUE
);

CREATE INDEX package_slug_index ON package USING HASH (slug);
CREATE INDEX package_code_index ON package USING HASH (code);
CREATE INDEX package_island_filter_index ON package USING HASH (island_filter);
