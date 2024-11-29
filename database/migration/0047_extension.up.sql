-- Extension Table
CREATE TABLE extension (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    name character varying(120) NOT NULL,
    slug character varying(120) NOT NULL,
    icon varchar(120) NOT NULL, -- the MDI icon for this page path
    data text NOT NULL DEFAULT '{}',
    is_active boolean NOT NULL DEFAULT true
);

CREATE INDEX extension_name_index ON extension USING HASH (name);
CREATE INDEX extension_slug_index ON extension USING HASH (slug);
