-- Component Table
CREATE TABLE component (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    name character varying(300) NOT NULL,
    icon character varying(100) NOT NULL, -- mdi icon name
    description character varying(3000) NOT NULL, -- documentation of the component
    class_name character varying(300) NOT NULL,    -- the component class name of the
    html_tag character varying(300) NOT NULL, -- the html tag of the component
    child_tag_constraints character varying(120)[] NOT NULL, -- the html tags that are allowed to be children of this component
    is_active boolean NOT NULL DEFAULT true
);

CREATE INDEX component_name_index ON component USING HASH (name);
CREATE INDEX component_class_name_index ON component USING HASH (class_name);
CREATE INDEX component_html_tag_index ON component USING HASH (html_tag);
