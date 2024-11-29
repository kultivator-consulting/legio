-- Location Table
CREATE TABLE location (
                          id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                          created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                          modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                          deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
                          name varchar(300) NOT NULL,
                          code varchar(300) NOT NULL,
                          ordering integer NOT NULL,
                          is_active boolean NOT NULL DEFAULT true
);

CREATE INDEX location_name_index ON location USING HASH (name);
CREATE INDEX location_code_index ON location USING HASH (code);
