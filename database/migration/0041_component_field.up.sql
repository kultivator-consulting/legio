-- Component Field Table
CREATE TABLE component_field (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    component_id uuid NOT NULL,
    name character varying(300) NOT NULL,
    description character varying(3000) NOT NULL, -- documentation of the component field
    data_type character varying(60) NOT NULL, -- the field data type, 'string', 'number', 'array', etc.
    editor_type character varying(60) NOT NULL, -- editor to be used to edit the field, 'text', 'number', 'date', etc.
    validation character varying(160) NOT NULL, -- validation rules for this field
    default_value character varying(300) NOT NULL, -- default value for this field
    is_active boolean NOT NULL DEFAULT true,
    CONSTRAINT component_field_component_id_fkey FOREIGN KEY (component_id) REFERENCES component
        ON DELETE CASCADE
);

CREATE INDEX component_field_component_id_index ON component_field USING HASH (component_id);
CREATE INDEX component_field_name_index ON component_field USING HASH (name);
