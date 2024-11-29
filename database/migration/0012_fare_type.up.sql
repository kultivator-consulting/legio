-- Fare Type Table
CREATE TABLE fare_type (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    name varchar(300) NOT NULL,
    ordering integer NOT NULL,
    description text NOT NULL,
    is_international boolean NOT NULL DEFAULT FALSE,
    is_active boolean NOT NULL DEFAULT true
);
