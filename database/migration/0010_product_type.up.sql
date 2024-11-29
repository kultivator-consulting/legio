-- Product Type Table
CREATE TABLE product_type (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    name varchar(300) NOT NULL,
    priority integer NOT NULL,
    slug varchar(300) NOT NULL,
    terms text NOT NULL,
    is_pass_type boolean NOT NULL DEFAULT false,
    is_pass_service boolean NOT NULL DEFAULT false,
    is_active boolean NOT NULL DEFAULT true
);

CREATE INDEX product_type_name_index ON product_type USING HASH (name);
CREATE INDEX product_type_slug_index ON product_type USING HASH (slug);
