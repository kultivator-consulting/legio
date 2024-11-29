-- Extra Table
CREATE TABLE extra (
                                  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
                                  created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
                                  deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
                                  name varchar(300) DEFAULT NULL,
                                  banner_image varchar(300) DEFAULT NULL,
                                  code varchar(150) DEFAULT NULL,
                                  description text,
                                  unit_price decimal DEFAULT 0.00,
                                  is_active boolean NOT NULL DEFAULT true
);

CREATE INDEX extras_name_index ON extra USING HASH (name);
CREATE INDEX extras_code_index ON extra USING HASH (code);
