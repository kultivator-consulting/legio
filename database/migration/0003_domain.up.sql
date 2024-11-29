-- Domain Table
CREATE TABLE domain (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    name varchar(300) NOT NULL,
    description varchar(2000) DEFAULT NULL,
    is_cms_controlled boolean NOT NULL DEFAULT false,
    is_under_maintenance boolean NOT NULL DEFAULT false,
    is_active boolean NOT NULL DEFAULT false
);

CREATE INDEX domain_name_index ON domain USING HASH (name);
