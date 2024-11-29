-- Pass Service Table
CREATE TABLE pass_service (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    created timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    modified timestamp with time zone DEFAULT (now() AT TIME ZONE 'UTC'),
    pass_id uuid NOT NULL,
    name varchar(300) NOT NULL,
    duration int NOT NULL,
    duration_type varchar(300) NOT NULL, -- hours, days, weeks, months
    CONSTRAINT pass_id_fkey FOREIGN KEY (pass_id) REFERENCES pass
        ON DELETE CASCADE
);

CREATE INDEX pass_service_pass_id_index ON pass_service USING HASH (pass_id);
CREATE INDEX pass_service_name_index ON pass_service USING HASH (name);

