-- Pass Table
CREATE TABLE pass (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    deleted timestamp with time zone DEFAULT ('infinity'::timestamp AT TIME ZONE 'UTC'),
    parent_pass_id uuid DEFAULT NULL,
    name varchar(300) NOT NULL,
    code varchar(120) NOT NULL,
    operator_codes varchar(30)[],
    duration int NOT NULL,
    duration_type varchar(300) NOT NULL, -- hours, days, weeks, months
    image varchar(400) NOT NULL,
    image_info varchar(600) NOT NULL,
    description text NOT NULL,
    is_popular boolean NOT NULL DEFAULT FALSE,
    is_top_up boolean NOT NULL DEFAULT FALSE,
    ordering int NOT NULL DEFAULT 0,
    is_active boolean NOT NULL DEFAULT TRUE
);

CREATE INDEX pass_parent_pass_id_index ON pass USING HASH (parent_pass_id);
CREATE INDEX pass_name_index ON pass USING HASH (name);
CREATE INDEX pass_code_index ON pass USING HASH (code);

